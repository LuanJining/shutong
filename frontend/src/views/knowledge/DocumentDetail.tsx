import api_frontend from "@/api/api_frontend"
import utils from "@/utils"
import { DownloadOutlined } from "@ant-design/icons"
import { Form, Tag, } from "antd"
import dayjs from 'dayjs'
import { useEffect, useState } from "react"
import { useLocation } from 'react-router-dom'
import FileUploader from "./Preview"
import "./styles/add-knowledge.scss"

export default function KonwledgeDetail() {
    const fromPage: any = useLocation().state
    const [info, setInfo] = useState<any>(null)
    const [classesName, setClassesName] = useState<string>('')

    useEffect(() => { fromPage?.documentId && getInfo() }, [fromPage?.documentId])
    useEffect(() => { info && getSpaces() }, [info])

    const getInfo = async () => {
        const r: any = await api_frontend.documentDetail(fromPage?.documentId)
        setInfo(r.data)
    }

    const getSpaces = async () => {
        const { data } = await api_frontend.getSpaces()
        data.map(({ name, id }: any) => ({ label: name, value: id }))
        const lv1 = getClass(data, info?.space_id, 'sub_spaces')
        const lv2 = getClass(lv1?.subArray, info?.sub_space_id, 'classes')
        const lv3 = getClass(lv2?.subArray, info?.class_id,)
        setClassesName(`${lv1?.name}/${lv2?.name}/${lv3?.name}`)
    }

    const getClass = (spaceArray: any[], spaceId: number, subKey?: string) => {
        const result: any = {
            name: '',
            subArray: []
        }
        spaceArray.map((v: any) => {
            if (v.id === spaceId) {
                subKey && (result.subArray = v[subKey])
                result.name = v.name
            }
        })
        return result
    }



    const BaseInfo = () => (<div className="form-box">
        <Form.Item label='创建者'>{info?.creator_nick_name}</Form.Item>

        <Form.Item label='所属部门'>
            {info?.department}
        </Form.Item>

        <Form.Item label='创建时间'>
            {dayjs(info?.created_at).format('YYYY-MM-DD HH:mm')}
        </Form.Item>

        <Form.Item label='所属分类'>
            <div className="flex flex-wrap al-start">
                {classesName}
            </div>
        </Form.Item>

        <Form.Item label='标签'>
            <div className="flex flex-wrap al-start">
                {info?.tags?.split('、').map((v: string) => <Tag color="blue" className="pdL12 pdR12 mgR16 mgB12" key={v}>{v}</Tag>)}
            </div>
        </Form.Item>

        <Form.Item label='摘要'>
            {info?.summary}
        </Form.Item>

        <Form.Item label='版本'>
            {info?.version}
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

                <div style={{
                    maxWidth: 'calc(100vw - 506px)'
                }} className="right-content flex1 flex flex-col">
                    <Form.Item style={{ marginBottom: 12 }} wrapperCol={{ offset: 1 }}>
                        <div className="flex al-center">
                            <div className="hg-fs elli fw-bold">{info?.title}</div>
                            {info?.id && <DownloadOutlined onClick={download} className="lg-fs mgL16 pointer fw-bold primary-blue" />}
                        </div>
                    </Form.Item>

                    <Form.Item className="h-100p" wrapperCol={{ offset: 1 }}>
                        <FileUploader
                            styles={{ maxHeight: 'calc(100vh - 190px)' }}
                            file={null} type={"url"} fileType={info?.file_type === '.pdf' ? 'pdf' : 'docx'} />
                    </Form.Item>
                </div>
            </Form>
        </div>
    )
}
