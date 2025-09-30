import ContructImg from '@/assets/images/construct.png'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'

export default function PageConstruct() {
    const navigate = useNavigate()
    return (
        <div className='h-100p'>
            <div onClick={() => { navigate(-1) }} className="back pointer pd24">
                <ArrowLeftOutlined className="ls-fs mgR6" />
                <span>返回</span>
            </div>
            <div className='flex-center flex-col'>
                <img style={{
                    width: '30%',
                }} src={ContructImg} alt="" />
                <div
                    className='fs36'
                    style={{ letterSpacing: 5, }}
                >网站建设中 . . .</div>
            </div>
        </div>
    )
}
