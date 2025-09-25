import React from 'react';
import { Button, Result } from 'antd';
import { useNavigate } from 'react-router-dom';
import "./index.scss"

const App: React.FC = () => {
    const navigate = useNavigate()
    return <Result
        status="404"
        title="404"
        subTitle="您访问的页面不存在"
        extra={<Button type="primary" onClick={() => navigate('/')}>返回首页</Button>}
        className='h-100p flex-center flex-col'
    />
}

export default App;