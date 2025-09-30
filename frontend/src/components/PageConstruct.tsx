import ContructImg from '@/assets/images/construct.png'

export default function PageConstruct() {
    return (
        <div className='flex-center flex-col h-100p'>
            <img style={{
                width: '30%',
            }} src={ContructImg} alt="" />
            <div
                className='fs36'
                style={{ letterSpacing: 5,}}
            >网站建设中 . . .</div>
        </div>
    )
}
