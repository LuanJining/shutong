import api_frontend from "@/api/api_frontend";
import { Par_Change_Pwd } from "@/types/api";
import utils from "@/utils";
import { message, Modal, } from "antd";
import { useForm } from "antd/es/form/Form";
import { Form, Input } from "antd/lib";
import { useNavigate } from "react-router-dom";
type modalType = {
    open: boolean;
    setOpen: (par: boolean) => void
}
export default function ResetPwdModal({ open, setOpen }: modalType) {
    const [form] = useForm()
    const navigate = useNavigate()

    const close = () => {
        setOpen(false)
        form.resetFields()
    };

    const onFinish = async (values: Par_Change_Pwd) => {
        utils.setLoading(true)
        await api_frontend.changePwd(values)
        message.success('修改成功')
        navigate('/login')
        close()
        utils.setLoading(false)
    }

    return (
        <Modal
            className="reset-pwd-modal"
            title="修改密码"
            open={open}
            onCancel={close}
            maskClosable={false}
            centered
            onOk={() => { form.submit() }}
        >
            <Form form={form}
                className="mgT32" onFinish={onFinish} labelCol={{ span: 4 }} labelAlign="right">
                <Form.Item label="原始密码" name="old_password" rules={[{ required: true, message: '请输入原始密码' }]}>
                    <Input.Password placeholder="请输入原始密码" />
                </Form.Item>
                <Form.Item label="新密码" name="new_password" rules={[{ required: true, message: '请输入新密码' }]}  >
                    <Input.Password placeholder="请输入新密码" />
                </Form.Item>

                <Form.Item dependencies={['new_password']} label="重复密码" name="confirmPwd" rules={[{
                    required: true, validator: (_rules, value: any) => {
                        console.log(value)
                        if (!value) return Promise.reject('请输入重复密码')
                        if (form.getFieldValue('new_password') !== value) return Promise.reject('新密码与重复密码不一致')
                        return Promise.resolve()
                    }
                }]} >
                    <Input.Password placeholder="输入新密码" />
                </Form.Item>
            </Form>

        </Modal>
    );
}
