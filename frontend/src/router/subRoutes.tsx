import { Spin } from "antd";
import { Suspense, lazy } from "react"
const Home = lazy(() => import("@/views/home/Index"));


/* 懒加载需要添加loading组件 */
const withLoadingComponent = (comp: JSX.Element) => (
    <Suspense
        fallback={
            <div className="flex-center" style={{ height: "calc(100vh - 65px)" }}>
                <Spin size="large" />
            </div>
        }
    >
        {comp}
    </Suspense>
);
export default [
    {
        path: "home",
        element: withLoadingComponent(<Home />),
        meta: {
            title: "首页",
        },
    },
    // {
    //     meta: { title: "账户管理", },
    //     children: [
    //         {
    //             path: "company",
    //             element: withLoadingComponent(<Company />),
    //             meta: { title: "组织管理", },
    //         },
    //         {
    //             path: "company-detail",
    //             element: withLoadingComponent(<CompanyDetail />),
    //             meta: {
    //                 title: "单位详情",
    //                 show: false
    //             },
    //         },
    //         {
    //             path: "company-create",
    //             element: withLoadingComponent(<CompanyCreate />),
    //             meta: {
    //                 title: "新增单位",
    //                 show: false
    //             },
    //         },
    //         {
    //             path: "member",
    //             element: withLoadingComponent(<Member />),
    //             meta: { title: "成员管理", },
    //         },
    //         {
    //             path: "role",
    //             element: withLoadingComponent(<Role />),
    //             meta: { title: "角色管理", },
    //         },
    //         {
    //             path: "role-info",
    //             element: withLoadingComponent(<RoleInfo />),
    //             meta: { title: "角色信息", show: false },
    //         },

    //     ]
    // },

]

