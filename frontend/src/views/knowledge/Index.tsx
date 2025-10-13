import "./index.scss"
import { useCallback, useEffect, useState } from 'react'
import _optsEnum from "@/config/optsEnum"
import _opts from '@/config/opts';
import api_frontend from "@/api/api_frontend"
import DocumentManagement from "./DocumentManagement"
import PageNav from "@/components/PageNav"
import { Props_Self_Nav } from "@/types/pages";
import SubPaceManagement from "./SubPaceManagement";
import { Button, Col, Row } from "antd";
import { useNavigate } from "react-router-dom";

export default function Index() {
    const navigate = useNavigate()
    const [pathKey, setPathKey] = useState<number>(-1)
    const [pageType, setPageType] = useState<'document' | 'management'>('document')
    const [items, setItems] = useState<(Props_Self_Nav & any)[]>([])

    useEffect(() => { setPathKey(items?.[0]?.key ?? '') }, [items])
    useEffect(() => { getSpaces() }, [])

    const getSpaces = useCallback(async () => {
        const { data: { spaces } }: any = await api_frontend.getSpaces()
        setItems(spaces.map((v: any) => ({ ...v, label: v.name, key: v.id })))
        setPathKey(spaces?.at(0)?.id)
    }, [pathKey])

    return (
        <div className='app-knowledge'>
            <PageNav pathKey={pathKey} pathArray={items} setPathKey={setPathKey} getSpaces={getSpaces} />
            <div className="nav-content flex1">
                <div className="flex al-center pdT24">
                    <div onClick={() => { setPageType(pageType === 'document' ? 'management' : 'document') }} className="hg-fs fw-bold pdL24 pointer">{pageType === 'document' ? '文档' : '子空间'}</div>
                    <i className="iconfont icon-exchange fs20 mgL24"></i>
                    <div onClick={() => { setPageType(pageType === 'document' ? 'management' : 'document') }} className="nm-fs fw-bold pdL24 pointer">{pageType === 'document' ? '子空间' : '文档'}</div>
                </div>
                <Row className="mgT24 pdL24 pdR24">
                    <div className="sm-fs primary-gray">知识库的所有文件都在这里显示，整个知识库都可以链接到应用引用或通过 Chat 插件进行索引。</div>
                    {
                        pageType === 'document' && <Col className="text-right" flex={1}>
                            <Button
                                style={{ width: 150 }}
                                onClick={() => { navigate('/knowledge/add') }}
                                type="primary">
                                + 添加文档
                            </Button>
                        </Col>
                    }
                </Row>
                {
                    pageType === 'document'
                        ? <DocumentManagement space_id={+pathKey} />
                        : <SubPaceManagement
                            space_id={+pathKey}
                            getSpaces={getSpaces} />
                }
            </div>
        </div>
    )
}
