import "../styles/layout-content.scss";
import Nav from "./Nav"
import { Outlet } from "react-router-dom";

export default function Content() {
    return (
        <div className="layout-content">
            <Nav />
            <div className="flex1 outlet">
                <Outlet />
            </div>
        </div>
    );
}
