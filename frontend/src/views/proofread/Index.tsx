import api_frontend from '@/api/api_frontend'
import caches from '@/config/caches'
import utils from '@/utils'
import storage from '@/utils/storage'
import { CheckCircleOutlined, CloseCircleOutlined, CloudUploadOutlined } from '@ant-design/icons'
import { Button, message } from 'antd'
import { useState } from 'react'
import DocumentUploader from './components/DocumentUploader'
import SuggestionPanel from './components/SuggestionPanel'
import './proofread.scss'

interface ReviewSuggestion {
    type: string
    severity: string
    position: number
    original_text: string
    suggested_text: string | null
    reason: string
    knowledge_source?: string | null
    knowledge_document_id?: number | null
    document_content?: string  // 文档内容（第一条建议携带）
}

export default function Proofread() {
    const [sessionId, setSessionId] = useState<string>('')
    const [fileName, setFileName] = useState<string>('')
    const [fileType, setFileType] = useState<string>('')
    const [suggestions, setSuggestions] = useState<ReviewSuggestion[]>([])
    const [isReviewing, setIsReviewing] = useState(false)
    const [uploadedFile, setUploadedFile] = useState<any>(null)
    const [documentContent, setDocumentContent] = useState<string>('')
    const [highlightedLine, setHighlightedLine] = useState<number | null>(null)
    
    // 上传文档
    const handleUpload = async (file: any) => {
        try {
            utils.setLoading(true)
            const formData = new FormData()
            formData.append('file', file)
            
            const response = await api_frontend.reviewUpload(formData)
            
            if (response.code === 200) {
                setSessionId(response.data)
                setFileName(file.name)
                // 获取文件类型
                const ext = file.name.substring(file.name.lastIndexOf('.'))
                setFileType(ext)
                setUploadedFile(file)
                
                // 读取文件内容用于左侧显示
                if (ext === '.txt' || ext === '.md') {
                    const reader = new FileReader()
                    reader.onload = (e) => {
                        const content = e.target?.result as string || ''
                        setDocumentContent(content)
                    }
                    reader.readAsText(file)
                } else {
                    // PDF/Word等待后端Pandoc转换后显示
                    setDocumentContent('')
                }
                
                message.success('文档上传成功')
            } else {
                message.error(response.message || '上传失败')
            }
        } catch (error: any) {
            message.error(error.message || '上传失败')
        } finally {
            utils.setLoading(false)
        }
        
        return false // 阻止默认上传行为
    }
    
    // 开始审查
    const startReview = async () => {
        if (!sessionId) {
            message.warning('请先上传文档')
            return
        }
        
        setIsReviewing(true)
        setSuggestions([])
        
        try {
            // 获取token
            const authInfo = storage.get(caches.AUTH_INFO)
            const token = authInfo?.access_token
            
            if (!token) {
                message.error('请先登录')
                setIsReviewing(false)
                return
            }
            
            // 构建SSE URL
            const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
            const url = `${baseURL}/review/${sessionId}/suggestions?fileName=${encodeURIComponent(fileName)}&fileType=${encodeURIComponent(fileType)}&checkFormat=true&verifyReferences=true&suggestContent=true`
            
            console.log('Starting review with URL:', url)
            console.log('Token:', token ? 'Present' : 'Missing')
            
            // 使用fetch来处理SSE（EventSource不支持自定义header）
            const response = await fetch(url, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'text/event-stream'
                }
            })
            
            console.log('Response status:', response.status)
            console.log('Response headers:', response.headers.get('content-type'))
            
            if (!response.ok) {
                const errorText = await response.text()
                console.error('Request failed:', errorText)
                throw new Error(`审查请求失败: ${response.status}`)
            }
            
            const reader = response.body?.getReader()
            const decoder = new TextDecoder()
            
            if (reader) {
                let buffer = ''
                let suggestionCount = 0
                
                while (true) {
                    const { done, value } = await reader.read()
                    
                    if (done) {
                        console.log('SSE stream ended')
                        setIsReviewing(false)
                        if (suggestionCount === 0) {
                            message.success('审查完成，未发现问题')
                        } else {
                            message.success(`审查完成，共发现 ${suggestionCount} 条建议`)
                        }
                        break
                    }
                    
                    buffer += decoder.decode(value, { stream: true })
                    const lines = buffer.split('\n')
                    buffer = lines.pop() || ''
                    
                    for (const line of lines) {
                        const trimmedLine = line.trim()
                        if (!trimmedLine) continue
                        
                        console.log('SSE line:', line)
                        
                        // Spring SSE格式：data:{...}
                        if (trimmedLine.startsWith('data:')) {
                            const data = trimmedLine.substring(5).trim()
                            console.log('SSE data:', data)
                            
                            if (data && data !== '') {
                                try {
                                    const suggestion = JSON.parse(data)
                                    console.log('Parsed suggestion:', suggestion)
                                    
                                    // 特殊处理：文档内容标记
                                    if (suggestion.type === 'DOCUMENT_CONTENT' && suggestion.document_content) {
                                        console.log('Received document content, length:', suggestion.document_content.length)
                                        setDocumentContent(suggestion.document_content)
                                        // 不计入建议数量，不显示在列表中
                                    } else {
                                        // 正常建议
                                        setSuggestions(prev => [...prev, suggestion])
                                        suggestionCount++
                                    }
                                } catch (e) {
                                    console.error('解析建议失败:', e, 'data:', data)
                                }
                            }
                        }
                        
                        // 检查结束标记：:DONE 或 event:
                        if (trimmedLine.startsWith(':DONE') || trimmedLine === ':' || trimmedLine.startsWith('event:')) {
                            console.log('Received end signal:', trimmedLine)
                        }
                    }
                }
            }
        } catch (error: any) {
            message.error(error.message || '审查失败')
            setIsReviewing(false)
        }
    }
    
    // 点击建议，跳转到对应位置
    const handlePositionClick = (position: number) => {
        setHighlightedLine(position)
        // 滚动到对应行
        const lineElement = document.getElementById(`line-${position}`)
        if (lineElement) {
            lineElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
        }
    }
    
    // 重置
    const handleReset = () => {
        setSessionId('')
        setFileName('')
        setFileType('')
        setSuggestions([])
        setIsReviewing(false)
        setUploadedFile(null)
        setDocumentContent('')
        setHighlightedLine(null)
    }
    
    return (
        <div className="proofread-container">
            {/* 头部 */}
            <div className="proofread-header">
                <div className="header-content">
                    <h1 className="title">智能校对</h1>
                    <p className="subtitle">上传文档，AI帮您检查格式、引用和内容</p>
                </div>
            </div>
            
            {/* 主体内容 */}
            <div className="proofread-body">
                {!uploadedFile ? (
                    // 上传区域
                    <DocumentUploader onUpload={handleUpload} />
                ) : (
                    // 审查区域
                    <div className="review-area">
                        {/* 文件信息卡片 */}
                        <div className="file-info-card">
                            <div className="file-info">
                                <CloudUploadOutlined className="file-icon" />
                                <div className="file-details">
                                    <div className="file-name">{fileName}</div>
                                    <div className="file-size">
                                        {(uploadedFile.size / 1024).toFixed(2)} KB
                                    </div>
                                </div>
                            </div>
                            <div className="file-actions">
                                {!isReviewing && suggestions.length === 0 && (
                                    <Button
                                        type="primary"
                                        size="large"
                                        onClick={startReview}
                                        icon={<CheckCircleOutlined />}
                                    >
                                        开始审查
                                    </Button>
                                )}
                                <Button
                                    size="large"
                                    onClick={handleReset}
                                    icon={<CloseCircleOutlined />}
                                >
                                    重新上传
                                </Button>
                            </div>
                        </div>
                        
                        {/* 左右分屏区域 */}
                        {(isReviewing || suggestions.length > 0) && (
                            <div className="split-view">
                                {/* 左侧：原文档 */}
                                <div className="left-pane">
                                    <div className="pane-header">
                                        <h3>原文档</h3>
                                    </div>
                                    <div className="pane-content">
                                        {documentContent ? (
                                            <div className="document-text">
                                                {documentContent.split('\n').map((line, index) => (
                                                    <div 
                                                        key={index}
                                                        id={`line-${index}`}
                                                        className={`document-line ${highlightedLine === index ? 'highlighted' : ''}`}
                                                    >
                                                        <span className="line-number">{index + 1}</span>
                                                        <span className="line-content">{line || ' '}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        ) : isReviewing ? (
                                            <div className="no-preview">
                                                <p>正在处理文档...</p>
                                                <p className="hint">PDF/Word正在转换为Markdown</p>
                                            </div>
                                        ) : (
                                            <div className="no-preview">
                                                <p>等待开始审查</p>
                                                <p className="hint">点击"开始审查"按钮</p>
                                            </div>
                                        )}
                                    </div>
                                </div>
                                
                                {/* 右侧：建议面板 */}
                                <div className="right-pane">
                                    <div className="pane-header">
                                        <h3>审查建议</h3>
                                    </div>
                                    <div className="pane-content">
                                        <SuggestionPanel
                                            suggestions={suggestions}
                                            isReviewing={isReviewing}
                                            onPositionClick={handlePositionClick}
                                        />
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    )
}

