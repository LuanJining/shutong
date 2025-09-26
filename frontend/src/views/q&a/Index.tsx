import "./index.scss"
import ServiceImg from "@/assets/images/service.png"
import { useNavigate } from "react-router-dom"
import { ArrowLeftOutlined } from "@ant-design/icons"
import { Input, Select } from "antd"

export default function Index() {
    const navigate = useNavigate()

    return (
        <div className='app-qa'>
            <div onClick={() => { navigate(-1) }} className="pointer">
                <ArrowLeftOutlined className="ls-fs mgR6" />
                <span>返回</span>
            </div>

            <div className="chat-content flex flex1 flex-col">

                <div className="flex1 chat-container flex-center">

                    <div className="chat-desc flex-center flex-col">
                        <img src={ServiceImg} alt="" />
                        <div className="ai-chart">AI 职工助手</div>
                        <div className="primary-gray mn-fs">我是AI 职工助手，你可以向我提出具体问题，我将学习平台内相关内容来为你解答～</div>
                    </div>

                </div>


                <div className="send-box flex flex-col">
                    <Input.TextArea style={{ resize: 'none' }} variant="borderless" rows={4} placeholder="请输入你的问题和需求" />
                    <div className="flex1 flex jf-end al-end">
                        <Select
                            style={{ width: 150 }}
                            value={1}
                            options={[
                                { label: '知识库', value: 1 }
                            ]} />
                    </div>
                </div>
            </div>
        </div>
    )
}
