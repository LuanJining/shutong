import "../styles/layout-nav.scss"
import { useLocation, useNavigate } from 'react-router-dom';
import subRoutes from '@/router/subRoutes'
import { ROUTESELF_TYPE } from "@/types/common";
import { useMemo } from "react";
import _caches from "@/config/caches";

export default function Nav() {
    const navigate = useNavigate()
    const pathname = useLocation().pathname

    const hasNav = useMemo(() => _caches.STYLE_WHITE_PATH.includes(pathname), [pathname])

    const handleClick = (e: any) =>  navigate(`${e.path}`)

    return <div
        style={{ display: hasNav ? 'flex' : 'none' }}
        className="layout-nav flex flex-col al-center">
        {
            subRoutes.filter((v: ROUTESELF_TYPE) => v?.meta?.show).map((v: ROUTESELF_TYPE) => <div
                key={v.path}
                onClick={() => { handleClick(v) }}
                className="nav-item mn-fs flex-center flex-col">
                <div className="nav-icon">{v.meta?.icon}</div>
                <span>{v.meta?.title}</span>
            </div>)
        }
    </div>
}


