import "./index.scss"
import { useState } from "react"
import { ArrowLeftOutlined, TeamOutlined, UserOutlined } from "@ant-design/icons"
import { useNavigate } from "react-router-dom"
import { Menu, MenuProps } from "antd"
import IconBook from '@/assets/icons/icon-book.png'
import Personal from "./Personal"
import RolesManament from "./RolesManament"
import UsersManagement from "./UsersManagement"

type MenuItem = Required<MenuProps>['items'][number];

const items: MenuItem[] = [
    {
        key: 'personal',
        label: '个人信息',
        icon: <UserOutlined style={{fontSize:18}}/>,
    },
    {
        key: 'users',
        label: '用户管理',
        icon:<TeamOutlined style={{fontSize:18}}/>,
    },
    {
        key: 'roles',
        label: '角色管理',
        icon: <img style={{ width: 20, height: 20, objectFit: 'cover' }} src={IconBook} />,
    },
];


const compEnum: any = {
    personal: <Personal />,
    users: <UsersManagement />,
    roles: <RolesManament />,
}

export default function Index() {
    const navigate = useNavigate()
    const [selectKeys, setKeys] = useState<string[]>(['personal'])

    const menuClick: MenuProps['onClick'] = (e) => {
        setKeys([e.key])
    };

    return (
        <div className='app-personal flex'>
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
