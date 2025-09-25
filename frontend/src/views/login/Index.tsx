import api_frontend from "@/api/api_frontend";
import "./index.scss"
import { Button, Form, Input, message } from "antd";
import storage from "@/utils/storage";
import _cache from "@/config/_cache";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { setIsLogin, setUserInfo } from "@/store/systemSlice";

export default function Index() {
    const navigate = useNavigate()
    const dispatch = useDispatch()

    const onFinish = async (values: any) => {
        const r: any = await api_frontend.login(values)
        storage.set(_cache.AUTH_INFO, r.data)
        dispatch(setIsLogin({ isLogin: true }))
        dispatch(setUserInfo({ userInfo: r.data.user }))
        message.success('登录成功')
        navigate('/home')
    }

    return (
        <div className="app-login flex-center flex-col">
            <div className="text-center mgB24 hg-fs fw-bold">用户登录</div>

            <Form labelCol={{ span: 5 }} className="login-form" onFinish={onFinish}>

                <Form.Item name='login' label='用户名'>
                    <Input placeholder="用户名/手机号/邮箱" />
                </Form.Item>

                <Form.Item name='password' label='密码'>
                    <Input.Password placeholder="请输入密码" />
                </Form.Item>

                <Form.Item className="flex-center mgT32">
                    <Button>取消</Button>
                    <Button className="mgL24" htmlType="submit" type="primary">确定</Button>
                </Form.Item>
            </Form>

        </div>
    )
}
