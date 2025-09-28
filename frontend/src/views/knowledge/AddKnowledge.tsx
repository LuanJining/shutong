import { Checkbox, DatePicker, DatePickerProps, Form, GetProps, Input, message, Radio, Select, Upload } from "antd"
import "./styles/add-knowledge.scss"
import Dragger from "antd/es/upload/Dragger"
import IconPdf from '@/assets/icons/icon-pdf.png'
import IconWendang from '@/assets/icons/icon-wendang.png'
import { CheckCircleFilled, CloudUploadOutlined, DeleteOutlined, EyeFilled } from "@ant-design/icons"
import CateSelect from "./CateSelect"
import { useEffect, useState } from "react"
import utils from "@/utils"
import FileUploader from "./Preview"
import { useSelector } from "react-redux"
import { useForm } from "antd/es/form/Form"
import _opts from '@/config/_opts';
import api_frontend from "@/api/api_frontend"
// type RangePickerProps = GetProps<typeof DatePicker.RangePicker>;

const initFileInfo: any = {
    file: null, fileType: ''
}

export default function AddKnowledge() {
    const [form] = useForm()
    const [cateOpen, setCopen] = useState<boolean>(false)
    const [isBase, setIsBase] = useState<boolean>(true)
    const [fileInfo, setFileInfo] = useState<any>(initFileInfo)
    const [spaces, setSpaces] = useState<any[]>([])
    const userInfo: any = useSelector((state: any) => state.systemSlice.userInfo)

    useEffect(() => { getSpaces() }, [])
    useEffect(() => { form.setFieldValue('department', userInfo?.department) }, [userInfo])

    const getSpaces = async () => {
        const { data: { spaces } }: any = await api_frontend.getSpaces()
        setSpaces(spaces.map((v: any) => ({ label: v.name, value: v.id })))
    }
    // const onOk = (value: DatePickerProps['value'] | RangePickerProps['value']) => {
    //     console.log('onOk: ', value);
    // };

    const beforeUpload = (file: any) => {
        const whiteArr = [
            'application/msword',
            'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
            'application/pdf',
        ]
        if (!whiteArr.includes(file.type)) {
            message.error("doc/docx/pdf", 1)
        }
        return whiteArr.includes(file.type) || Upload.LIST_IGNORE
    }

    const customRequest = async ({ file, onSuccess, onError }: any) => {
        try {
            utils.setLoading(true)
            console.log(file)
            setFileInfo({
                file: file,
                fileType: [
                    'application/msword',
                    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
                ].includes(file.type)
                    ? 'docx'
                    : 'pdf'
            })
            form.setFieldValue('file_name', file.name)
            // const formatData = utils.getFormData('file', file)
            // await api_upload.dataFlowUpload(formatData, curFlowerId)
            message.success("上传成功", 1)
            onSuccess()
            utils.setLoading(false)
        } catch {
            onError()
        }
    }

    const BaseInfo = () => (<div className="form-box">
        <Form.Item
            valuePropName="fileList"
            getValueFromEvent={utils.normFile}
            name='bigFile'
            rules={[{ required: true, message: '请上传文档' }]}
        >
            <Dragger
                maxCount={1}
                showUploadList={false}
                customRequest={customRequest}
                beforeUpload={beforeUpload}
                onRemove={() => {
                    setFileInfo(initFileInfo)
                    form.setFieldValue('file_name', '')
                    return true
                }}
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

        <Form.Item name='space_id' rules={[{ required: true, message: '请选择知识库' }]} label='知识库'>
            <Select options={spaces} />
        </Form.Item>

        <Form.Item label='创建者'>{userInfo?.nickname}</Form.Item>

        <Form.Item name='department' label='所属部门'>
            {/* <Input placeholder="所属部门" /> */}
            {userInfo?.department}
        </Form.Item>

        {/* <Form.Item name='' label='知识分类'>
            <Input />
        </Form.Item> */}

        <Form.Item rules={[{ required: true, message: '请选择标签' }]} name='tags' label='知识标签'>
            <Select mode="multiple" options={[
                { label: '知识', value: '知识' },
                { label: '文档', value: '文档' },
                { label: '管理', value: '管理' },
            ]} />
        </Form.Item>

        {/* <Form.Item label='下架时间'>
            <DatePicker
                style={{ width: '100%' }}
                showTime
                onChange={(value, dateString) => {
                    console.log('Selected  Time: ', value);
                    console.log('Formatted Selected Time: ', dateString);
                }}
                onOk={onOk}
            />
        </Form.Item> */}

        {/* <Form.Item name='role' label='权限'>
            <Radio.Group
                value={1}
                options={[
                    { value: 1, label: '可查看' },
                    { value: 2, label: '可应用' },
                ]}
            />
        </Form.Item> */}

        {/* <Form.Item style={{ marginTop: -24 }} wrapperCol={{ offset: 5 }} >
            <Select options={[]} />
        </Form.Item> */}

        <Form.Item initialValue={true} name='need_approval' label='审批'>
            <Radio.Group
                disabled
                options={[
                    { value: true, label: '审批' },
                    { value: false, label: '无需审批' },
                ]}
            />
        </Form.Item>

        {/* <Form.Item noStyle shouldUpdate={(pre: any, cur: any) => pre.approve !== cur.approve}>
            {
                ({ getFieldValue }) => {
                    const isApprove: boolean = getFieldValue('approve') === 1
                    return isApprove ? <Form.Item style={{ marginTop: -24 }} wrapperCol={{ offset: 5 }} >
                        <Select options={[]} />
                    </Form.Item> : <></>
                }
            }
        </Form.Item>

        <Form.Item style={{ marginBottom: 0 }} label='版本'>
            <Input />
        </Form.Item>  */}

    </div>)


    const Approve = () => (<div className="form-box">

        <Form.Item label='提交人'>
            {/* <Input /> */}
            {userInfo?.nickname}
        </Form.Item>

        <Form.Item label='所属部门'>
            {/* <Input /> */}
            {userInfo?.department}
        </Form.Item>

        {/* <Form.Item label='下一节点审批人'>
            <Input />
        </Form.Item> */}

        <Form.Item label='是否最终节点'>
            是
        </Form.Item>

        <Form.Item initialValue={_opts.URGENCY[0].value} name='urgency' label='紧急程度'>
            <Select options={_opts.URGENCY} />
        </Form.Item>

        <Form.Item name='remark' label='备注'>
            <Input.TextArea rows={4} style={{ resize: 'none' }} />
        </Form.Item>

        <Form.Item>
            <Checkbox defaultChecked />
            <span className="mgL6">流程结束后通知我</span>
        </Form.Item>
    </div>)

    const onFinish = async () => {
        utils.setLoading(true)
        const values: any = form.getFieldsValue([
            'urgency', 'tags', 'summary', 'need_approval', 'space_id', 'file_name'
        ])
        const par: any = utils.getFormData({
            ...values,
            tags: values.tags.join('、'),
            visibility: 'public',
            department: userInfo?.department,
            created_by: userInfo?.id,
            file: fileInfo.file
        })
        await api_frontend.uploadFile(par)
        message.success('提交成功')
        utils.setLoading(false)
    }

    const onFinishFailed = (errorInfo: any) => {
        console.log(errorInfo)
    }

    return (
        <div className='add-knowledge h-100p'>
            <Form
                form={form}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                requiredMark={false}
                className="h-100p flex al-stretch" labelAlign="left" colon={false} labelCol={{ span: isBase ? 5 : 7 }}>

                <div className="left-content">
                    <div className="top-box flex al-center space-between">
                        <div className="flex al-center chose">
                            <div onClick={() => { setIsBase(true) }} className="chose-txt">基本信息</div>
                            <div className="mgL16 mgR16">|</div>
                            <div onClick={() => {
                                form.validateFields(['space_id', 'tags']).then(() => {
                                    setIsBase(false)
                                })
                            }} className="chose-txt">流程审批</div>
                            <div
                                style={{ left: `${isBase ? 0 : 92}px` }}
                                className="active-line"></div>
                        </div>

                        <div onClick={() => { form.submit() }} className="btn-confirm">提交</div>
                    </div>
                    {isBase ? <BaseInfo /> : <Approve />}
                </div>

                <div className="right-content flex1 flex flex-col">
                    <div className="file-info flex al-center space-between">
                        <div className="flex al-center">
                            {
                                fileInfo.fileType ?
                                    <div className="icon-img mgR16"> <img src={fileInfo.fileType === 'pdf' ? IconPdf : IconWendang} alt="" /></div>
                                    : <></>
                            }
                            <div className="primary-gray nm-fs elli">{form.getFieldValue('file_name')}</div>
                        </div>

                        <div className="flex flex1 jf-end al-center">

                            <div className="flex al-center white-nowrap" style={{ marginRight: 100 }}>
                                <CheckCircleFilled style={{ color: '#52CC6F' }} />
                                <span className="mgL12">处理完成</span>
                            </div>

                            <div className="flex al-center">
                                <EyeFilled className="mgR12" />
                                <DeleteOutlined onClick={() => {
                                    setFileInfo(initFileInfo)
                                    form.setFieldValue('file_name', '')
                                }} className="pointer" />
                            </div>
                        </div>
                    </div>

                    <Form.Item name='file_name' className="mgT24" labelCol={{ span: 1 }} label='名称' rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>

                    <Form.Item rules={[{ required: true, }]} name='summary' labelCol={{ span: 1 }} label='摘要'>
                        <Input.TextArea rows={4} style={{ resize: 'none' }} />
                    </Form.Item>

                    <Form.Item className="h-100p" wrapperCol={{ offset: 1 }}>
                        <FileUploader
                            file={fileInfo.file}
                            type='file'
                            fileType={fileInfo.fileType}
                            styles={{
                                maxHeight: 'calc(100vh - 380px)',
                                maxWidth: ' calc(100vw - 680px)'
                            }}
                        />
                    </Form.Item>

                </div>
            </Form>

            <CateSelect open={cateOpen} setOpen={setCopen} />
        </div>
    )
}
