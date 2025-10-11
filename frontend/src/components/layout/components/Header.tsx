import "../styles/layout-header.scss";
import storage from "@/utils/storage";
import IconUser from "@/assets/icons/icon-user.png"
import { Layout, } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { BellOutlined, LogoutOutlined } from "@ant-design/icons";
import { useMemo } from "react";
import BannerImg from '@/assets/images/banner.png'
import _caches from "@/config/caches";

const { Header } = Layout;

export default function LayoutHeader() {
    const navigate = useNavigate()
    const pathname = useLocation().pathname

    const hasBgImg = useMemo(() => !_caches.STYLE_WHITE_PATH.includes(pathname), [pathname])

    const logout = () => {
        storage.clear()
        navigate("/login")
    }

    return (
        <Header
            style={{ backgroundImage: hasBgImg ? `url(${BannerImg})` : '' }}
            className="layout-header al-center jf-end">

            <BellOutlined onClick={() => {
                navigate('/notification')
            }} className="hg-fs mgR24 pointer" />

            <img onClick={() => {
                navigate('/personal')
            }} className="user-img pointer" src={IconUser} alt="" />

            <LogoutOutlined
                title="登出"
                onClick={logout}
                className="pointer"
            />
        </Header>
    );
}
