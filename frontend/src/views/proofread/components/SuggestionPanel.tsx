import { BulbOutlined, CheckCircleOutlined, ExclamationCircleOutlined, InfoCircleOutlined } from '@ant-design/icons'
import { Empty, Segmented, Spin } from 'antd'
import { useMemo, useState } from 'react'
import SuggestionCard from './SuggestionCard'

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

interface SuggestionPanelProps {
    suggestions: ReviewSuggestion[]
    isReviewing: boolean
    onPositionClick?: (position: number) => void
}

export default function SuggestionPanel({ suggestions, isReviewing, onPositionClick }: SuggestionPanelProps) {
    const [filterType, setFilterType] = useState<string>('all')
    
    // 统计
    const stats = useMemo(() => {
        const errorCount = suggestions.filter(s => s.severity === 'ERROR').length
        const warningCount = suggestions.filter(s => s.severity === 'WARNING').length
        const infoCount = suggestions.filter(s => s.severity === 'INFO').length
        
        return {
            total: suggestions.length,
            error: errorCount,
            warning: warningCount,
            info: infoCount
        }
    }, [suggestions])
    
    // 过滤建议
    const filteredSuggestions = useMemo(() => {
        if (filterType === 'all') return suggestions
        return suggestions.filter(s => s.severity === filterType)
    }, [suggestions, filterType])
    
    // 图标映射
    const severityIcons = {
        'ERROR': <ExclamationCircleOutlined className="severity-icon error" />,
        'WARNING': <InfoCircleOutlined className="severity-icon warning" />,
        'INFO': <BulbOutlined className="severity-icon info" />
    }
    
    return (
        <div className="suggestion-panel">
            {/* 统计卡片 - 精简版 */}
            <div className="stats-compact">
                <div className="stats-row">
                    <span className="stats-title">审查结果</span>
                    {isReviewing && <Spin size="small" />}
                    <div className="stats-badges">
                        <span className="stat-badge">总数 {stats.total}</span>
                        {stats.error > 0 && <span className="stat-badge error">错误 {stats.error}</span>}
                        {stats.warning > 0 && <span className="stat-badge warning">警告 {stats.warning}</span>}
                        {stats.info > 0 && <span className="stat-badge info">提示 {stats.info}</span>}
                    </div>
                </div>
            </div>
            
            {/* 筛选器 */}
            {suggestions.length > 0 && (
                <div className="filter-bar">
                    <Segmented
                        value={filterType}
                        onChange={(value) => setFilterType(value as string)}
                        options={[
                            {
                                label: `全部 ${stats.total}`,
                                value: 'all'
                            },
                            {
                                label: `错误 ${stats.error}`,
                                value: 'ERROR',
                                disabled: stats.error === 0
                            },
                            {
                                label: `警告 ${stats.warning}`,
                                value: 'WARNING',
                                disabled: stats.warning === 0
                            },
                            {
                                label: `提示 ${stats.info}`,
                                value: 'INFO',
                                disabled: stats.info === 0
                            }
                        ]}
                    />
                </div>
            )}
            
            {/* 建议列表 */}
            <div className="suggestions-list">
                {isReviewing && suggestions.length === 0 ? (
                    <div className="reviewing-placeholder">
                        <Spin size="large" />
                        <p>正在审查文档，请稍候...</p>
                    </div>
                ) : suggestions.length === 0 ? (
                    <Empty
                        description="未发现任何问题"
                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                    >
                        <CheckCircleOutlined style={{ fontSize: 48, color: '#52c41a' }} />
                        <p style={{ marginTop: 16, color: '#52c41a', fontSize: 16, fontWeight: 500 }}>
                            文档格式规范，无需修改
                        </p>
                    </Empty>
                ) : (
                    filteredSuggestions.map((suggestion, index) => (
                        <SuggestionCard
                            key={index}
                            suggestion={suggestion}
                            icon={severityIcons[suggestion.severity as keyof typeof severityIcons]}
                            onPositionClick={onPositionClick}
                        />
                    ))
                )}
            </div>
        </div>
    )
}

