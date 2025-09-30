import { ROUTESELF_TYPE } from "@/types/common";
import { Spin } from "antd";
import { Suspense, lazy } from "react"

import IconPen from '@/assets/icons/icon-pen.png'
import IconMask from '@/assets/icons/icon-mask.png'
import IconAns from '@/assets/icons/icon-ans.png'
import IconBook from '@/assets/icons/icon-book.png'
import IconSearch from '@/assets/icons/icon-search.png'

const PageConstruct = lazy(() => import("@/components/PageConstruct"));

const Qa = lazy(() => import("@/views/q&a/Index"));
const Home = lazy(() => import("@/views/home/Index"));
const AddKnowledge = lazy(() => import("@/views/knowledge/AddKnowledge"));
const Notification = lazy(() => import("@/views/notification/Index"));
const Personal = lazy(() => import("@/views/personal/Index"));
const Knowledge = lazy(() => import("@/views/knowledge/Index"));
const DocumentDetail = lazy(() => import("@/views/knowledge/DocumentDetail"));

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

const routes: ROUTESELF_TYPE[] = [
    {
        path: "home",
        element: withLoadingComponent(<Home />),
        meta: {
            title: "首页",
            show: false
        },
    },
    {
        path: "write",
        element: withLoadingComponent(<PageConstruct />),
        meta: {
            title: "写作",
            icon: <img src={IconPen} />,
            show: true
        },
    },
    {
        path: "approve",
        element: withLoadingComponent(<PageConstruct />),
        meta: {
            title: "审核",
            icon: <img src={IconMask} />,
            show: true
        },
    },
    {
        path: "q&a",
        element: withLoadingComponent(<Qa />),
        meta: {
            title: "问答",
            icon: <img src={IconAns} />,
            show: true
        },
    },
    {
        path: "proofread",
        element: withLoadingComponent(<PageConstruct />),
        meta: {
            title: "校对",
            icon: <img src={IconSearch} />,
            show: true
        },
    },

    {
        path: "knowledge",
        element: withLoadingComponent(<Knowledge />),
        meta: {
            title: "知识库",
            icon: <img src={IconBook} />,
            show: true
        },
    },
    {
        path: "document/detail",
        element: withLoadingComponent(<DocumentDetail />),
        meta: {
            title: "文档详情",
            show: false
        },
    },

    {
        path: "knowledge/add",
        element: withLoadingComponent(<AddKnowledge />),
        meta: {
            title: "新增文档知识",
            show: false
        },
    },

    {
        path: "notification",
        element: withLoadingComponent(<Notification />),
        meta: {
            title: "消息通知",
            show: false
        },
    },

    {
        path: "personal",
        element: withLoadingComponent(<Personal />),
        meta: {
            title: "个人中心",
            show: false
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

export default routes

