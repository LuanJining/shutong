import { Checkbox, DatePicker, DatePickerProps, Form, GetProps, Input, message, Radio, Select, Upload } from "antd"
import "./styles/add-knowledge.scss"
import Dragger from "antd/es/upload/Dragger"
import IconPdf from '@/assets/icons/icon-pdf.png'
import { CheckCircleFilled, CloudUploadOutlined, DeleteOutlined, EyeFilled } from "@ant-design/icons"
import CateSelect from "./CateSelect"
import { useState } from "react"
import utils from "@/utils"
type RangePickerProps = GetProps<typeof DatePicker.RangePicker>;

export default function AddKnowledge() {
    const [cateOpen, setCopen] = useState<boolean>(false)
    const [isBase, setIsBase] = useState<boolean>(true)

    const onOk = (value: DatePickerProps['value'] | RangePickerProps['value']) => {
        console.log('onOk: ', value);
    };

    const beforeUpload = (file: any) => {
        console.log(file)
        const whiteArr = [
            'application/msword',
            'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
            'application/pdf',
        ]
        if (!whiteArr.includes(file.type)) {
            message.error("doc/docx/pdf", 1)
        }
        return whiteArr.includes(file.type) ? true : Upload.LIST_IGNORE
    }

    const customRequest = async (options: any) => {
        try {
            utils.setLoading(true)
            const formatData = utils.getFormData('file', options.file)
            // await api_upload.dataFlowUpload(formatData, curFlowerId)
            message.success("上传成功", 1)
            options.onSuccess()
            utils.setLoading(false)
        } catch {
            options.onError()
        }
    }

    const BaseInfo = () => (<div className="form-box">
        <Form.Item
            valuePropName="fileList"
            getValueFromEvent={utils.normFile}
            name='bigFile'
        >
            <Dragger
                customRequest={customRequest}
                beforeUpload={beforeUpload}
            >
                <p className="ant-upload-drag-icon">
                    <CloudUploadOutlined />
                </p>
                <p className="ant-upload-text">上传文件，请<span className="primary-blue pointer">点击上传</span></p>
                <p className="ant-upload-hint">
                    仅支持 docx/doc/PDF 文件，文件大小不超过50M
                </p>
            </Dragger>
        </Form.Item>
        <Form.Item label='创建者'>杨柳</Form.Item>

        <Form.Item label='所属部门'>
            <Select options={[]} />
        </Form.Item>

        <Form.Item label='知识分类'>
            <Input />
        </Form.Item>

        <Form.Item label='知识标签'>
            <Input />
        </Form.Item>

        <Form.Item label='下架时间'>
            <DatePicker
                style={{ width: '100%' }}
                showTime
                onChange={(value, dateString) => {
                    console.log('Selected Time: ', value);
                    console.log('Formatted Selected Time: ', dateString);
                }}
                onOk={onOk}
            />
        </Form.Item>

        <Form.Item name='role' label='权限'>
            <Radio.Group
                value={1}
                options={[
                    { value: 1, label: '可查看' },
                    { value: 2, label: '可应用' },
                ]}
            />
        </Form.Item>

        <Form.Item style={{ marginTop: -24 }} wrapperCol={{ offset: 5 }} >
            <Select options={[]} />
        </Form.Item>

        <Form.Item name='approve' label='审批'>
            <Radio.Group
                value={1}
                options={[
                    { value: 1, label: '审批' },
                    { value: 2, label: '无需审批' },
                ]}
            />
        </Form.Item>
        <Form.Item noStyle shouldUpdate={(pre: any, cur: any) => pre.approve !== cur.approve}>
            {
                ({ getFieldValue }) => {
                    const isApprove: boolean = getFieldValue('approve') === 1
                    return isApprove ? <Form.Item style={{ marginTop: -24 }} wrapperCol={{ offset: 5 }} >
                        <Select options={[]} />
                    </Form.Item> : <></>
                }
            }
        </Form.Item>

        <Form.Item label='版本'>
            <Input />
        </Form.Item>

    </div>)


    const Approve = () => (<div className="form-box">

        <Form.Item label='提交人'>
            <Input />
        </Form.Item>

        <Form.Item label='所属部门'>
            <Input />
        </Form.Item>

        <Form.Item label='下一节点审批人'>
            <Input />
        </Form.Item>

        <Form.Item label='是否最终节点'>
            <Input />
        </Form.Item>

        <Form.Item label='紧急程度'>
            <Input />
        </Form.Item>

        <Form.Item label='备注'>
            <Input />
        </Form.Item>

        <Form.Item>
            <Checkbox />
            <span className="mgL6">流程结束后通知我</span>
        </Form.Item>
    </div>)

    return (
        <div className='add-knowledge h-100p'>
            <Form className="h-100p flex al-stretch" labelAlign="left" colon={false} labelCol={{ span: isBase ? 5 : 7 }}>

                <div className="left-content">
                    <div className="top-box flex al-center space-between">
                        <div className="flex al-center chose">
                            <div onClick={() => { setIsBase(true) }} className="chose-txt">基本信息</div>
                            <div className="mgL16 mgR16">|</div>
                            <div onClick={() => { setIsBase(false) }} className="chose-txt">流程审批</div>
                            <div
                                style={{ left: `${isBase ? 0 : 92}px` }}
                                className="active-line"></div>
                        </div>

                        <div className="btn-confirm">提交</div>
                    </div>
                    {isBase ? <BaseInfo /> : <Approve />}
                </div>

                <div className="right-content flex1 flex flex-col">
                    <div className="file-info flex al-center space-between">
                        <div className="flex al-center">
                            <div className="icon-img">
                                <img src={IconPdf} alt="" />
                            </div>
                            <div className="primary-gray sm-fs elli">关于印发《安全生产文明施工“党政同责”暂行规定》的通知（核西南建[2016]41号）.pdf</div>
                        </div>

                        <div className="flex al-center">

                            <div className="flex al-center white-nowrap" style={{ marginRight: 100 }}>
                                <CheckCircleFilled style={{ color: '#52CC6F' }} />
                                <span className="mgL12">处理完成</span>
                            </div>

                            <div className="flex al-center">
                                <EyeFilled className="mgR12" />
                                <DeleteOutlined />
                            </div>
                        </div>
                    </div>

                    <Form.Item className="mgT24" labelCol={{ span: 1 }} label='名称' name='name' rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>

                    <Form.Item labelCol={{ span: 1 }} label='摘要'>
                        <Input.TextArea rows={4} style={{ resize: 'none' }} />
                    </Form.Item>

                    <Form.Item className="h-100p" wrapperCol={{ offset: 1 }}>
                        <div className="pdf-container">
                            <object
                                type="application/pdf"
                                data={''}
                                width="100%"
                                height='100%'
                            >
                                <div className="p-center flex-center">
                                    <a href={''}>您的浏览器不支持 PDF 文件，请下载后查看。</a>
                                </div>
                            </object>
                        </div>
                    </Form.Item>

                </div>
            </Form>

            <CateSelect open={cateOpen} setOpen={setCopen} />
        </div>
    )
}
