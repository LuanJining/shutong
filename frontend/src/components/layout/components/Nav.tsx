import "../styles/layout-nav.scss"
import { useNavigate } from 'react-router-dom';
import subRoutes from '@/router/subRoutes'
import { ROUTESELF_TYPE } from "@/types/common";

export default function Nav() {
    const navigate = useNavigate()

    const handleClick = (e: any) => navigate(`${e.key}`)

    return <div className="layout-nav flex flex-col al-center">
        {
            subRoutes.filter((v: ROUTESELF_TYPE) => v?.meta?.show).map((v: ROUTESELF_TYPE) => <div
                key={v.path}
                className="nav-item mn-fs flex-center flex-col">
                <div className="nav-icon">{v.meta?.icon}</div>
                <span>{v.meta?.title}</span>
            </div>)
        }
    </div>
}


