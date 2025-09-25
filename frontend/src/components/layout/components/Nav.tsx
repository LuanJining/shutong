import "../styles/layout-nav.scss"
import {  useNavigate } from 'react-router-dom';

export default function Nav() {
    const navigate = useNavigate()

    const handleClick = (e: any) => navigate(`${e.key}`)

    return <div className="layout-nav"></div>
}


