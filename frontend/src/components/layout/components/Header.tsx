import "../styles/layout-header.scss";
import storage from "@/utils/storage";
import IconUser from "@/assets/icons/icon-user.png"
import { Layout, } from "antd";
import { useNavigate } from "react-router-dom";
import { LogoutOutlined, UserOutlined } from "@ant-design/icons";

const { Header } = Layout;

export default function LayoutHeader() {
    const navigate = useNavigate()

    const logout = () => {
        storage.clear()
        navigate("/login")
    }

    return (
        <Header className="layout-header al-center jf-end">

            <img className="user-img" src={IconUser} alt="" />

            <LogoutOutlined
                title="登出"
                onClick={logout}
                className="pointer"
            />
        </Header>
    );
}
