import api_frontend from '@/api/api_frontend'
import caches from '@/config/caches'
import utils from '@/utils'
import storage from '@/utils/storage'
import { CheckCircleOutlined, CloseCircleOutlined, CloudUploadOutlined, DownloadOutlined } from '@ant-design/icons'
import { Button, message } from 'antd'
import { useRef, useState } from 'react'
import DocumentUploader from './components/DocumentUploader'
import SuggestionPanel from './components/SuggestionPanel'
import './proofread.scss'

interface ReviewSuggestion {
    id?: string  // 建议ID
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
    const [acceptedSuggestions, setAcceptedSuggestions] = useState<string[]>([])

    // 使用ref来存储关键状态，防止重新渲染时丢失
    const sessionIdRef = useRef<string>('')
    const fileNameRef = useRef<string>('')
    const fileTypeRef = useRef<string>('')

    // 调试：跟踪状态变化
    console.log('Proofread component render - sessionId:', sessionId, 'fileName:', fileName, 'fileType:', fileType)

    // 上传文档
    const handleUpload = async (file: any) => {
        try {
            utils.setLoading(true)
            const formData = new FormData()
            formData.append('file', file)

            const response = await api_frontend.reviewUpload(formData)

            if (response.code === 200) {
                console.log('Upload response data:', response.data)
                setSessionId(response.data)
                setFileName(file.name)
                // 获取文件类型
                const lastDotIndex = file.name.lastIndexOf('.')
                const ext = lastDotIndex > 0 ? file.name.substring(lastDotIndex) : '.txt'
                setFileType(ext)
                setUploadedFile(file)

                // 同时更新ref
                sessionIdRef.current = response.data
                fileNameRef.current = file.name
                fileTypeRef.current = ext

                console.log('SessionId set to:', response.data)

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
                                        console.log('Adding suggestion with ID:', suggestion.id)
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

    // 接受建议
    const handleAcceptSuggestions = async (suggestionIds: string[]) => {
        console.log('handleAcceptSuggestions called with:', suggestionIds)
        console.log('Current sessionId:', sessionId)
        console.log('Current fileName:', fileName)
        console.log('Current fileType:', fileType)

        // 使用ref中的值作为备用
        const currentSessionId = sessionId || sessionIdRef.current
        const currentFileName = fileName || fileNameRef.current
        const currentFileType = fileType || fileTypeRef.current

        console.log('Using values - sessionId:', currentSessionId, 'fileName:', currentFileName, 'fileType:', currentFileType)

        if (!currentSessionId || currentSessionId.trim() === '') {
            console.error('SessionId is empty or null:', currentSessionId)
            message.warning('会话ID缺失，请重新上传文档')
            return
        }

        if (!currentFileName) {
            message.warning('文件名信息缺失，请重新上传文档')
            return
        }

        if (!currentFileType) {
            message.warning('文件类型信息缺失，请重新上传文档')
            return
        }

        try {
            utils.setLoading(true)

            const requestData = {
                sessionId: currentSessionId,
                fileName: currentFileName,
                fileType: currentFileType,
                acceptedSuggestionIds: suggestionIds,
                applyAll: false
            }


            console.log('Accept suggestions request data:', requestData)

            const response = await api_frontend.acceptSuggestions(requestData)

            if (response.code === 200) {
                setAcceptedSuggestions(prev => [...prev, ...suggestionIds])
                message.success('建议已成功应用')
            } else {
                message.error(response.message || '应用建议失败')
            }
        } catch (error: any) {
            message.error(error.message || '应用建议失败')
        } finally {
            utils.setLoading(false)
        }
    }

    // 接受所有建议
    const handleAcceptAllSuggestions = async () => {
        if (!sessionId) {
            message.warning('请先上传文档')
            return
        }

        try {
            utils.setLoading(true)

            const response = await api_frontend.acceptSuggestions({
                sessionId,
                fileName,
                fileType,
                applyAll: true
            })

            if (response.code === 200) {
                const allIds = suggestions.filter(s => s.id).map(s => s.id!)
                setAcceptedSuggestions(allIds)
                message.success('所有建议已成功应用')
            } else {
                message.error(response.message || '应用建议失败')
            }
        } catch (error: any) {
            message.error(error.message || '应用建议失败')
        } finally {
            utils.setLoading(false)
        }
    }

    // 下载修改后的文档
    const handleDownload = async () => {
        // 使用ref中的值作为备用
        const currentSessionId = sessionId || sessionIdRef.current
        const currentFileName = fileName || fileNameRef.current
        const currentFileType = fileType || fileTypeRef.current

        if (!currentSessionId) {
            message.warning('请先上传文档')
            return
        }

        try {
            utils.setLoading(true)

            const response = await api_frontend.downloadDocument({
                sessionId: currentSessionId,
                fileName: currentFileName,
                fileType: currentFileType
            })

            // 创建下载链接
            const blob = new Blob([response], { type: 'application/octet-stream' })
            const url = window.URL.createObjectURL(blob)
            const link = document.createElement('a')
            link.href = url
            link.download = fileName
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            window.URL.revokeObjectURL(url)

            message.success('文档下载成功')
        } catch (error: any) {
            message.error(error.message || '下载失败')
        } finally {
            utils.setLoading(false)
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
        setAcceptedSuggestions([])

        // 同时重置ref
        sessionIdRef.current = ''
        fileNameRef.current = ''
        fileTypeRef.current = ''
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
                                {suggestions.length > 0 && (
                                    <Button
                                        size="large"
                                        onClick={handleDownload}
                                        icon={<DownloadOutlined />}
                                        type="default"
                                    >
                                        下载文档
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
                                            onAcceptSuggestions={handleAcceptSuggestions}
                                            onAcceptAllSuggestions={handleAcceptAllSuggestions}
                                            acceptedSuggestions={acceptedSuggestions}
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

