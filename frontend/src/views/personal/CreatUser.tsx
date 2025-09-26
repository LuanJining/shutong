import api_frontend from "@/api/api_frontend";
import { Par_Users } from "@/types/api";
import utils from "@/utils";
import { Form, Input, message, Modal } from "antd";
import { useForm } from "antd/es/form/Form";

interface IProps {
    open: boolean, setOpen: Function, callback: Function
}


export default function CreatUser({ open, setOpen, callback }: IProps) {
    const [form] = useForm()

    const onFinish = async (values: Par_Users) => {
        try {
            utils.setLoading(true)
            await api_frontend.createUser(values)
            message.success('创建成功')
            utils.setLoading(false)
            setOpen(false)
            callback()
        } catch (e){
            utils.setLoading(false)
            throw(e)
        }
    }

    return (
        <Modal
            title='新增用户'
            centered
            open={open}
            destroyOnClose
            onCancel={() => { setOpen(false) }}
            onOk={() => { form.submit() }}
        >
            <Form
                form={form}
                onFinish={onFinish}
                labelCol={{ span: 4 }}
                wrapperCol={{ offset: 1 }}
                autoComplete="off"
                className="mgT24"
            >
                <Form.Item
                    label="用户名"
                    name="username"
                    rules={[{ required: true, message: '请输入用户名!' }]}
                >
                    <Input placeholder="请输入用户名" />
                </Form.Item>

                <Form.Item
                    label="昵称"
                    name="nickname"
                    rules={[{ required: true, message: '请输入昵称!' }]}
                >
                    <Input placeholder="请输入昵称" />
                </Form.Item>

                <Form.Item
                    label="手机号"
                    name="phone"
                    rules={[
                        { required: true, message: '请输入手机号!' },
                        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式!' },
                    ]}
                >
                    <Input placeholder="请输入手机号" />
                </Form.Item>

                <Form.Item
                    label="邮箱"
                    name="email"
                    rules={[
                        { required: true, message: '请输入邮箱!' },
                        { type: 'email', message: '请输入正确的邮箱格式!' },
                    ]}
                >
                    <Input placeholder="请输入邮箱" />
                </Form.Item>

                <Form.Item
                    label="密码"
                    name="password"
                    rules={[
                        { required: true, message: '请输入密码!' },
                        { min: 6, message: '密码至少6位!' },
                    ]}
                >
                    <Input.Password placeholder="请输入密码" />
                </Form.Item>

                <Form.Item
                    label="部门"
                    name="department"
                    rules={[{ required: true, message: '请输入部门!' }]}
                >
                    <Input placeholder="请输入部门" />
                </Form.Item>

                <Form.Item
                    label="公司"
                    name="company"
                    rules={[{ required: true, message: '请输入公司!' }]}
                >
                    <Input placeholder="请输入公司" />
                </Form.Item>
            </Form>
        </Modal>
    )
}
