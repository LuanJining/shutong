import "./styles/add-knowledge.scss"
import dayjs from 'dayjs'
import { Form, Tag, } from "antd"
import { useEffect, useState } from "react"
import FileUploader from "./Preview"
import _opts from '@/config/_opts';
import api_frontend from "@/api/api_frontend"
import { useLocation } from 'react-router-dom';
import { DownloadOutlined } from "@ant-design/icons"
import utils from "@/utils"

export default function KonwledgeDetail() {
    const fromPage: any = useLocation().state
    const [info, setInfo] = useState<any>(null)

    useEffect(() => { fromPage?.documentId && getInfo() }, [fromPage?.documentId])

    const getInfo = async () => {
        const r: any = await api_frontend.documentDetail(fromPage?.documentId)
        setInfo(r)
    }

    const BaseInfo = () => (<div className="form-box">
        <Form.Item label='创建者'>{info?.creator?.nickname}</Form.Item>

        <Form.Item label='所属部门'>
            {info?.department}
        </Form.Item>

        <Form.Item label='创建时间'>
            {dayjs(info?.created_at).format('YYYY-MM-DD HH:mm')}
        </Form.Item>

        <Form.Item label='标签'>
            <div className="flex flex-wrap al-start">
                {info?.tags?.split('、').map((v: string) => <Tag color="blue" className="pdL12 pdR12 mgR16 mgB12" key={v}>{v}</Tag>)}
            </div>
        </Form.Item>

        <Form.Item label='摘要'>
            {info?.summary}
        </Form.Item>
    </div>)

    const download = async () => {
        utils.setLoading(true)
        const res: any = await api_frontend.getFile(fromPage?.documentId)
        utils.downloadFromFlow(res, info?.file_name)
        utils.setLoading(false)
    }

    return (
        <div className='add-knowledge h-100p'>
            <Form
                requiredMark={false}
                className="h-100p flex al-stretch" labelAlign="left" colon={false} labelCol={{ span: 5 }}>

                <div className="left-content">
                    <div className="top-box flex al-center space-between">
                        基本信息
                    </div>
                    <BaseInfo />
                </div>

                <div className="right-content flex1 flex flex-col">
                    <Form.Item style={{ marginBottom: 12 }} wrapperCol={{ offset: 1 }}>
                        <div className="flex al-center">
                            <div className="hg-fs elli fw-bold">{info?.title}</div>
                            <DownloadOutlined onClick={download} className="lg-fs mgL16 pointer fw-bold primary-blue" />
                        </div>
                    </Form.Item>

                    <Form.Item className="h-100p" wrapperCol={{ offset: 1 }}>
                        <FileUploader
                            styles={{
                                maxHeight: 'calc(100vh - 190px)',
                            }}
                            file={null} type={"url"} fileType={info?.file_type === '.pdf' ? 'pdf' : 'docx'} />
                    </Form.Item>
                </div>
            </Form>
        </div>
    )
}
