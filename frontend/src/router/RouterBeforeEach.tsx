import router from "./index";
import _utils from "@/utils";
import storage from "@/utils/storage";
import { useEffect } from "react";
import { useSelector } from "react-redux";
import { useLocation, useNavigate, useRoutes } from "react-router-dom";

const whitelist = ["/login"];

export default function RouterBeforeEach() {
    const location = useLocation();
    const navigate = useNavigate();

    const outlet = useRoutes(router);
    const pathname = location.pathname;
    const isLogin = useSelector((state: any) => state.systemSlice.isLogin)

    useEffect(() => {
        const isRouteInWhitelist = whitelist.includes(location.pathname);
        if (!isRouteInWhitelist && !isLogin) {
            return navigate("/login");
        }
        if (pathname === "/login") {
            storage.clear()
        }
    }, [pathname, isLogin]);

    return outlet;
}
