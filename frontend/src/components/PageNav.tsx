import "./styles/page-nav.scss"
import { useNavigate } from "react-router-dom"
import {  Props_Self_Nav } from "@/types/pages";
import { ArrowLeftOutlined, BorderOuterOutlined, DeleteOutlined, EditOutlined, PlusCircleOutlined } from "@ant-design/icons"
import PageKnowledgeModal from "./PageKnowledgeModal";
import { useState } from "react";
import { message, Popconfirm } from "antd";
import utils from "@/utils";
import api_frontend from "@/api/api_frontend";

interface IProps {
    pathKey: number;
    setPathKey: (par: number) => void;
    pathArray: Props_Self_Nav[];
    getSpaces: () => void
}

export default function PageNav({ pathKey, pathArray, setPathKey, getSpaces }: IProps) {
    const navigate = useNavigate()
    const [modalType, setMType] = useState<'add' | 'update' | ''>('')

    const onDelete = async (space_id: string) => {
        try {
            utils.setLoading(true)
            await api_frontend.deleteSpace(space_id)
            message.success('删除成功')
            getSpaces()
            utils.setLoading(false)
        } catch (e) {
            utils.setLoading(false)
            throw (e)
        }
    }

    return (
        <div className='page-nav'>
            <div onClick={() => { navigate(-1) }} className="pointer pd24">
                <ArrowLeftOutlined className="sm-fs mgR6" />
                <span className="sm-fs">返回</span>
            </div>
            <div className="flex al-center space-between pdR12">
                <div className="fw-bold pdL24 mgB12 sm-fs">知识库</div>
                <PlusCircleOutlined
                    onClick={() => { setMType('add') }}
                    className="pointer fs20 primary-blue" />
            </div>
            <div className="nav-content">
                {
                    pathArray.map(({ key, label }: Props_Self_Nav & any) => {
                        return (<div
                            key={key}
                            onClick={() => { setPathKey(+key) }}
                            className={`nav-item flex space-between sm-fs ${+pathKey === +key ? 'nav-item-active' : ''}`}>
                            <span>{label}</span>
                            <div className="flex al-ceter icons-box">
                                <BorderOuterOutlined
                                    onClick={(e) => {
                                        e.stopPropagation()
                                        setPathKey(+key )
                                    }}
                                    title="管理" className="pointer nav-icon" />
                                <EditOutlined title='编辑' onClick={(e) => { e.stopPropagation(); setMType('update') }} className="pointer mgL12 nav-icon" />
                                <Popconfirm
                                    title="删除知识库"
                                    description="确认要删除知识库吗？"
                                    okText="确认"
                                    cancelText="取消"
                                    onPopupClick={(e) => { e.stopPropagation() }}
                                    onConfirm={(e) => { e?.stopPropagation(); onDelete(key) }}
                                >
                                    <DeleteOutlined title='删除' onClick={(e) => { e.stopPropagation() }} className="mgL12 pointer nav-icon" />
                                </Popconfirm>
                            </div>
                        </div>)
                    })
                }
            </div>
            <PageKnowledgeModal getSpaces={getSpaces} open={Boolean(modalType)} setOpen={() => { setMType('') }} item={
                modalType === 'update' ? pathArray.find(({ key }: Props_Self_Nav) => +key === +pathKey) : null
            } />
        </div>
    )
}
