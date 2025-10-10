import "./index.scss"
import api_frontend from "@/api/api_frontend";
import { Button, Checkbox, Form, Input, message } from "antd";
import storage from "@/utils/storage";
import _cache from "@/config/_caches";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { setIsLogin, setUserInfo } from "@/store/systemSlice";
import utils from "@/utils";

export default function Index() {
    const navigate = useNavigate()
    const dispatch = useDispatch()

    const onFinish = async (values: any) => {
        utils.setLoading(true)
        const r: any = await api_frontend.login(values)
        storage.set(_cache.AUTH_INFO, r.data)
        dispatch(setIsLogin({ isLogin: true }))
        dispatch(setUserInfo({ userInfo: r.data.user }))
        message.success('登录成功')
        navigate('/home')
        utils.setLoading(false)
    }

    return (
        <div className="app-login flex-center">
            <Form
                layout="vertical"
                className="login-form" onFinish={onFinish}>

                <div className="text-center mgB24 hg-fs fw-bold">用户登录</div>

                <Form.Item name='login' label='用户名'>
                    <Input placeholder="用户名/手机号/邮箱" />
                </Form.Item>

                <Form.Item name='password' label='密码'>
                    <Input.Password placeholder="请输入密码" />
                </Form.Item>

                <Form.Item>
                    <Checkbox />
                    <span className="mgL12">记住密码</span>
                </Form.Item>

                <Form.Item className="mgT32">
                    <Button style={{ width: '100%', height: 36 }} htmlType="submit" type="primary">登录</Button>
                </Form.Item>
            </Form>

        </div>
    )
}
