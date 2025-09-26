import { ArrowLeftOutlined } from "@ant-design/icons"
import "./index.scss"
import { useNavigate } from "react-router-dom"
import { Menu, MenuProps } from "antd"
import IconBook from '@/assets/icons/icon-book.png'
import IconNotifications from "@/assets/icons/icon-notifications.png"
import { useState } from "react"
import WaitDeal from "./WaitDeal"
import DoneDeal from "./DoneDeal"

type MenuItem = Required<MenuProps>['items'][number];

const items: MenuItem[] = [
    {
        key: 'knowledge',
        label: '知识管理',
        icon: <img style={{ width: 20, height: 20, objectFit: 'cover' }} src={IconBook} />,
        children: [
            {
                key: 'wait',
                label: '待处理',
            },
            {
                key: 'done',
                label: '已处理',
            },
        ],
    },
    {
        key: 'system',
        label: '系统通知',
        icon: <img style={{ width: 20, height: 20, objectFit: 'cover' }} src={IconNotifications} />,
    },
];

const compEnum: any = {
    wait: <WaitDeal />,
    done: <DoneDeal />,
    system: <WaitDeal />,
}

export default function Index() {
    const navigate = useNavigate()
    const [selectKeys, setKeys] = useState<string[]>(['wait'])

    const menuClick: MenuProps['onClick'] = (e) => {
        setKeys([e.key])
    };

    return (
        <div className='app-notification flex'>
            <div className="notification-nav">
                <div onClick={() => { navigate(-1) }} className="back pointer pd24">
                    <ArrowLeftOutlined className="ls-fs mgR6" />
                    <span>返回</span>
                </div>

                <Menu
                    onClick={menuClick}
                    style={{ width: '100%', background: 'transparent', border: 0 }}
                    defaultSelectedKeys={selectKeys}
                    defaultOpenKeys={['knowledge']}
                    mode="inline"
                    items={items}
                />
            </div>

            <div className="nav-content flex1">
                {compEnum[selectKeys[0]]}
            </div>
        </div>
    )
}
