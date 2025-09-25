import "../styles/layout-header.scss";
import storage from "@/utils/storage";
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
        <Header className="layout-header jf-end">
            <UserOutlined
                className="pointer lg-fs mgR24"
            />
            <LogoutOutlined
                title="登出"
                onClick={logout}
                className="pointer"
            />
        </Header>
    );
}
