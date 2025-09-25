import { Suspense, lazy } from "react";

import { Navigate } from "react-router-dom";
import { Spin } from "antd";
import { ROUTESELF_TYPE } from "@/types/common";
import subRoutes from "./subRoutes";

const Layout = lazy(() => import("@/components/layout/Index"));
const NotFound = lazy(() => import("@/views/not-found/Index"));
const Login = lazy(() => import("@/views/user/Index"));

/* 懒加载需要添加loading组件 */
const withLoadingComponent = (comp: JSX.Element) => (
    <Suspense
        fallback={
            <div className="flex-center" style={{ height: '100vh' }}>
                <Spin size="large" />
            </div>
        }
    >
        {comp}
    </Suspense>
);

const routes: ROUTESELF_TYPE[] = [
    {
        path: "/",
        element: <Navigate to="/home" />,
    },
    {
        path: "notFound",
        element: withLoadingComponent(<NotFound />),
    },
    {
        path: "*",
        element: withLoadingComponent(<NotFound />),
    },
    {
        path: "/login",
        element: withLoadingComponent(<Login />),
    },
    {
        element: withLoadingComponent(<Layout />),
        path: 'home',
        children: subRoutes,
    },
];

export default routes;
