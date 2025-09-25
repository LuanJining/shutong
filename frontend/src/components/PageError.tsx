import { Button,  Result } from 'antd'
import { useNavigate } from 'react-router-dom'

export default function Index() {
    const navigate = useNavigate()
    return (
        <div className='w-100p h-100p flex-center'>
            <Result
                status="500"
                title="网络异常"
                subTitle="网络异常，刷新重试下吧！"
                extra={<Button type="primary" onClick={()=>navigate(0)}>刷新</Button>}
            />
        </div>
    )
}
