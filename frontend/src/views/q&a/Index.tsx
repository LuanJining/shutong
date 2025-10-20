import api_frontend from "@/api/api_frontend"
import ServiceImg from "@/assets/images/service.png"
import _caches from "@/config/caches"
import { OPTIONS_TYPE } from "@/types/common"
import storage from "@/utils/storage"
import { LoadingOutlined } from "@ant-design/icons"
import { Input, message, Modal, Select } from "antd"
import { marked } from "marked"
import { useEffect, useRef, useState } from "react"
import "./index.scss"

export default function Index() {
    const [konwledges, setKonwledges] = useState<OPTIONS_TYPE[]>([])
    const [space_id, setSpaceId] = useState<string>('')
    const [loading, setloading] = useState<boolean>(false)
    const [words, setWords] = useState<string>('')
    const [isInit, setInit] = useState<boolean>(true)
    const abortController = useRef<any>(null)
    const chatBoxRef = useRef<any>(null)
    const chatContainerRef = useRef<any>(null)
    const fetchUrl: string = `${import.meta.env['VITE_API_URL']}/documents/chat/stream`;

    useEffect(() => { getSpaces() }, [])
    useEffect(() => { setSpaceId(konwledges?.[0]?.value) }, [konwledges])

    const getSpaces = async () => {
        const { data: { spaces } } = await api_frontend.getSpaces()
        setKonwledges(spaces.map(({ name, id }: any) => ({ label: name, value: id })))
    }

    useEffect(() => {
        if (loading) {
            setTimeout(scrollToBottom, 0);
        }
        return () => { Modal.destroyAll() }
    }, [loading])


    const scrollToBottom = () => {
        const chatContainer = chatContainerRef.current;
        if (chatContainer) {
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }
    };

    const sendMessage = () => {
        const messages = words.trim();
        const chatBox = chatBoxRef.current

        if (messages && !loading) {

            setInit(false)

            chatBox.innerHTML += `<div class="user-message"> <span>${messages}</span></div> `;
            setWords('')
            setloading(true)
            fetchLogStream(messages);
            scrollToBottom()
        } else if (loading) {
            // 中止请求
            if (abortController.current) {
                abortController.current?.abort();
            }
            setloading(false)
            abortController.current = null
        }
    }

    const fetchLogStream = async (keyword: string) => {
        const chatBox = chatBoxRef.current
        try {
            abortController.current = new AbortController();
            const token = storage.get(_caches.AUTH_INFO)?.access_token
            const response: any = await fetch(fetchUrl, {
                method: 'POST',
                signal: abortController.current?.signal,
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'text/event-stream',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    space_id,
                    question: keyword
                }),
            });
            if (!response.ok) throw new Error('请求失败');

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let accumulatedText = '';
            let messageDiv = null;

            messageDiv = document.createElement('div');
            messageDiv.className = 'ai-message';
            chatBox.appendChild(messageDiv);

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);

                const lines = chunk.split('\n');

                let currentEvent = '';
                for (const rawLine of lines) {
                    const line = rawLine.trim();
                    if (!line) continue;

                    if (line.startsWith('event:')) {
                        currentEvent = line.replace('event:', '').trim();
                        continue;
                    } else if (line.startsWith('data:') && currentEvent) {
                        const dataStr = line.replace('data:', '').trim();
                        if (!dataStr) continue;
                        try {
                            const data = JSON.parse(dataStr);
                            if (currentEvent === 'token') {
                                const content = data.content;
                                if (content) {
                                    accumulatedText += content;
                                    messageDiv.innerHTML = marked.parse(accumulatedText).toString();
                                    scrollToBottom()
                                }
                            } else if (currentEvent === 'done') {
                                console.log('done')
                            } else if (currentEvent === 'sources') {
                                const sources = data; // 假如是数组或对象
                                console.log('来源信息:', sources);
                            }
                            currentEvent = ''
                        } catch (e) {
                            console.warn('解析 event/data 失败:', e, '数据:', line);
                        }
                    }
                }
            }

        } catch (error: any) {
            if (error.name === 'AbortError') {
                chatBox.innerHTML += `<div class="ai-message">已停止生成。</div>`;
            } else {
                chatBox.innerHTML += `<div class="ai-message">发生错误，请重试。</div>`;
                console.error('日志流错误:', error);
            }
        } finally {
            setloading(false)
            abortController.current = null
        }
    }

    const handleKeyDown = (e: any) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            if (e.altKey) {
                const { selectionStart, selectionEnd, value } = e.target;
                const newValue =
                    value.substring(0, selectionStart) +
                    '\n' +
                    value.substring(selectionEnd);

                setWords(newValue);
                setTimeout(() => {
                    const textarea = e.target;
                    textarea.selectionStart = textarea.selectionEnd = selectionStart + 1;
                }, 0);
            } else {
                if (!loading) {
                    sendMessage();
                } else {
                    message.info('ai正在生成中，请稍后操作', 1.5)
                }
            }
        }
    };

    return (
        <div className='app-qa flex flex-col h-100p'>

            <div
                className={`flex1 chat-container flex-center ${!isInit ? 'display-none' : ''}`}>
                {
                    <div className="chat-desc flex-center flex-col">
                        <img src={ServiceImg} alt="" />
                        <div className="ai-chart">AI 职工助手</div>
                        <div className="primary-gray mn-fs">我是AI 职工助手，你可以向我提出具体问题，我将学习平台内相关内容来为你解答～</div>
                    </div>
                }
            </div>

            <div
                ref={chatContainerRef}
                className={`flex jf-center chat-box flex1 ${isInit ? 'display-none' : ''}`}>
                <div className="chat-container-box">
                    <div ref={chatBoxRef} className="chat-container pdB12">
                    </div>
                    {
                        loading && <div className="ai-message pdB12">
                            <span>正在思考...</span>
                            <LoadingOutlined className="fs16 mgL6" />
                        </div>
                    }
                </div>
            </div>

            <div className="flex jf-center">
                <div className="send-box flex flex-col">
                    <Input.TextArea
                        rows={4}
                        autoFocus
                        value={words}
                        variant="borderless"
                        onKeyDown={handleKeyDown}
                        style={{ resize: 'none' }}
                        onChange={(e: any) => { setWords(e.target.value) }}
                        placeholder="请输入你的问题和需求" />
                    <div className="flex1 flex jf-end al-end">
                        <Select
                            style={{ width: 150 }}
                            value={konwledges?.[0]?.value ?? -1}
                            options={konwledges} />

                        <div
                            onClick={() => { sendMessage() }}
                            className={`send-button ${loading ? 'loading-style' : 'normal-style'}`}
                        >
                            {
                                loading
                                    ? <svg id="stop-icon" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"
                                        xmlns="http://www.w3.org/2000/svg">
                                        <path fillRule="evenodd" clipRule="evenodd"
                                            d="M4.5 7.5a3 3 0 0 1 3-3h9a3 3 0 0 1 3 3v9a3 3 0 0 1-3 3h-9a3 3 0 0 1-3-3v-9Z" />
                                    </svg>
                                    : <svg id="send-icon" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"
                                        xmlns="http://www.w3.org/2000/svg">
                                        <path
                                            d="M3.478 2.404a.75.75 0 0 0-.926.941l2.432 7.905H13.5a.75.75 0 0 1 0 1.5H4.984l-2.432 7.905a.75.75 0 0 0 .926.94 60.519 60.519 0 0 0 18.445-8.986.75.75 0 0 0 0-1.218A60.517 60.517 0 0 0 3.478 2.404Z" />
                                    </svg>
                            }
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}
