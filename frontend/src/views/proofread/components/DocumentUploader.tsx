import { CloudUploadOutlined } from '@ant-design/icons'
import { Upload } from 'antd'

const { Dragger } = Upload

interface DocumentUploaderProps {
    onUpload: (file: any) => boolean | Promise<boolean>
}

export default function DocumentUploader({ onUpload }: DocumentUploaderProps) {
    const beforeUpload = (file: any) => {
        const whiteArr = [
            'application/pdf',
            'text/plain',
            'application/msword',
            'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        ]
        
        if (!whiteArr.includes(file.type)) {
            return Upload.LIST_IGNORE
        }
        
        // 文件大小限制 50MB
        const isLt50M = file.size / 1024 / 1024 < 50
        if (!isLt50M) {
            return Upload.LIST_IGNORE
        }
        
        return true
    }
    
    return (
        <div className="upload-area">
            <Dragger
                name="file"
                multiple={false}
                maxCount={1}
                showUploadList={false}
                beforeUpload={beforeUpload}
                customRequest={({ file }) => onUpload(file)}
                className="document-dragger"
            >
                <div className="upload-content">
                    <div className="upload-icon-wrapper">
                        <CloudUploadOutlined className="upload-icon" />
                    </div>
                    <h3 className="upload-title">
                        点击或拖拽文件到此区域上传
                    </h3>
                    <p className="upload-hint">
                        支持 PDF、Word、TXT 格式，文件大小不超过 50MB
                    </p>
                </div>
            </Dragger>
        </div>
    )
}

