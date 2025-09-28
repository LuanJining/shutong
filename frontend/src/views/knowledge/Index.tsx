import "./index.scss"
import { ArrowLeftOutlined } from "@ant-design/icons"
import { Menu, MenuProps, } from "antd"
import { useEffect, useState } from 'react'
import { useNavigate } from "react-router-dom"
import _optsEnum from "@/config/_optsEnum"
import _opts from '@/config/_opts';
import { OPTIONS_TYPE } from "@/types/common"
import KonwledgeManagement from "./KonwledgeManagement"
import api_frontend from "@/api/api_frontend"
import DocumentManagement from "./DocumentManagement"

type MenuItem = Required<MenuProps>['items'][number];

export default function Index() {
    const navigate = useNavigate()
    const [selectKeys, setKeys] = useState<string[]>(['konwledge'])
    const [items, setItems] = useState<MenuItem[]>([
        {
            key: "konwledge",
            label: '知识库管理',
        },
    ])

    useEffect(() => { getSpaces() }, [])

    const getSpaces = async () => {
        try {
            const { data: { spaces } }: any = await api_frontend.getSpaces()
            setItems([
                {
                    key: "konwledge",
                    label: '知识库管理',
                },
                {
                    key: "document",
                    label: '文档管理',
                    children: spaces.map((v: any) => ({
                        label: v.name,
                        key: v.id,
                    }))
                }
            ])
        } catch (e) {
            throw (e)
        }
    }

    const menuClick: MenuProps['onClick'] = (e) => {
        setKeys(e.keyPath)
    };

    const getComp = () => {
        if (selectKeys[0] === 'konwledge') return <KonwledgeManagement />
        if (selectKeys.includes('document')) {
            return <DocumentManagement space_id={selectKeys.at(0)!} />
        }
    }

    return (
        <div className='app-knowledge'>
            <div className="knowledge-nav">
                <div onClick={() => { navigate(-1) }} className="pointer pd24">
                    <ArrowLeftOutlined className="nm-fs mgR6" />
                    <span className="sm-fs">返回</span>
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
                {getComp()}
            </div>
        </div>
    )
}
