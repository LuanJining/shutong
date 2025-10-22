import { BookOutlined } from '@ant-design/icons'
import { Button, Tag } from 'antd'
import { ReactNode } from 'react'

interface ReviewSuggestion {
    type: string
    severity: string
    position: number
    original_text: string
    suggested_text: string | null
    reason: string
    knowledge_source?: string | null
    knowledge_document_id?: number | null
}

interface SuggestionCardProps {
    suggestion: ReviewSuggestion
    icon?: ReactNode
    onPositionClick?: (position: number) => void
}

export default function SuggestionCard({ suggestion, icon, onPositionClick }: SuggestionCardProps) {
    const { type, severity, original_text, suggested_text, reason, knowledge_source, position } = suggestion
    
    // 类型标签颜色
    const typeColors: Record<string, string> = {
        'FORMAT_ERROR': 'red',
        'PUNCTUATION': 'orange',
        'REFERENCE_OUTDATED': 'volcano',
        'REFERENCE_MISSING': 'red',
        'CONTENT_ENHANCEMENT': 'blue',
        'NUMBERING_ERROR': 'orange',
        'DATE_FORMAT': 'gold'
    }
    
    // 类型中文名
    const typeNames: Record<string, string> = {
        'FORMAT_ERROR': '格式错误',
        'PUNCTUATION': '标点符号',
        'REFERENCE_OUTDATED': '引用过期',
        'REFERENCE_MISSING': '缺少引用',
        'CONTENT_ENHANCEMENT': '内容建议',
        'NUMBERING_ERROR': '编号错误',
        'DATE_FORMAT': '日期格式'
    }
    
    const severityClass = severity.toLowerCase()
    
    return (
        <div 
            className={`suggestion-card ${severityClass}`}
            onClick={() => onPositionClick && onPositionClick(position)}
        >
            <div className="card-header">
                <div className="card-title">
                    {icon}
                    <span className="title-text">{reason}</span>
                </div>
                <Tag color={typeColors[type] || 'default'}>
                    {typeNames[type] || type}
                </Tag>
            </div>
            
            <div className="card-content">
                {original_text && (
                    <div className="text-block original">
                        <div className="block-label">原文 · 第{position + 1}行</div>
                        <div className="block-text">{original_text}</div>
                    </div>
                )}
                
                {suggested_text && (
                    <div className="text-block suggested">
                        <div className="block-label">建议</div>
                        <div className="block-text">{suggested_text}</div>
                    </div>
                )}
            </div>
            
            {knowledge_source && (
                <div className="card-footer">
                    <div className="knowledge-source">
                        <BookOutlined className="source-icon" />
                        <span className="source-text">来源：{knowledge_source}</span>
                    </div>
                    <Button type="link" size="small">
                        查看详情
                    </Button>
                </div>
            )}
        </div>
    )
}

