import "./index.scss"
import BannerImg from "@/assets/images/banner.png"
import LogoImg from "@/assets/images/logo.png"
import IconMouse from "@/assets/icons/icon-mouse.png"
import IconFileAdd from "@/assets/icons/icon-file-add.png"
import IconFile from "@/assets/icons/icon-file.png"
import IconTarget from "@/assets/icons/icon-target.png"
import IconFloder from "@/assets/icons/icon-floder.png"
import IconResume from "@/assets/icons/icon-resume.png"
import IconFloderAdd from "@/assets/icons/icon-floder-add.png"
import IconRight from "@/assets/icons/icon-right.png"

import { SearchOutlined } from "@ant-design/icons"
import { Col, Empty, Input, Row } from "antd"
import { useNavigate } from "react-router-dom"
import React, { useEffect, useMemo, useState } from "react"
import api_frontend from "@/api/api_frontend"

const countArray: any[] = [
    { label: '制度汇总', desc: '公司各部门制定指引', icon: <img src={IconFile} />, },
    { label: '法律法规', desc: '法律法规资源库', icon: <img src={IconResume} />, },
    { label: '项目经验', desc: '企业项目案例', icon: <img src={IconTarget} />, },
    { label: '行政后勤', desc: '后勤职能部门资源库', icon: <img src={IconFloder} />, },
    { label: '产品研发', desc: '企业产品档案库', icon: <img src={IconMouse} />, },
]

export default function Index() {
    const navigate = useNavigate()
    const [spaces, setSpaces] = useState<any[]>([])
    const [par, setPar] = useState<any>({})
    const [keywords, setKeywords] = useState<string>('')
    const [documents, setDocuments] = useState<any[]>([])
    const [isEmpty, setEmpty] = useState<boolean>(false)

    useEffect(() => { getHomePage() }, [])

    const getHomePage = async () => {
        const { data: { spaces } }: any = await api_frontend.homePage()
        setSpaces(spaces)
        const tempPar: any = {}
        spaces.map(({ id, sub_spaces }: any) => {
            tempPar[id] = sub_spaces?.at(0)?.id
        })
        setPar(tempPar)
    }

    const search = async () => {

        if (!keywords.trim()) {
            setDocuments([])
            setEmpty(false)
            return
        }

        const params: any = {
            limit: 5,
            query: keywords.trim()
        }
        const { data: { items } }: any = await api_frontend.search(params)
        setDocuments(items)
        setEmpty(items?.length === 0)
    }

    const firstOne: any = useMemo(() => spaces?.at(0), [spaces])
    const lastItems: any = useMemo(() => spaces?.slice(1, 5), [spaces])

    const getDocument = (spaceId: number) => {
        return spaces.find(({ id }: any) => id === spaceId)?.sub_spaces?.find(({ id }: any) => id === par[spaceId])?.documents ?? []
    }

    return (
        <div className="app-home">
            <div className="banner-box">
                <img src={BannerImg} alt="" />

                <div className="text-box flex al-center">
                    <div className="logo-img">
                        <img src={LogoImg} alt="" />
                    </div>
                    <div className="text-title">核工业西南建设集团有限公司 · 知识库</div>
                </div>

                <div className="search-box flex al-center">
                    <Input
                        value={keywords}
                        onPressEnter={search}
                        onChange={(e: any) => { setKeywords(e.target.value) }}
                        variant="borderless" placeholder="请输入搜索内容" />
                    <SearchOutlined onClick={search} className="hg-fs" />
                </div>
            </div>
            {
                !isEmpty && documents.length === 0 && <React.Fragment>

                    <div className="count-box flex al-center space-between">
                        {
                            countArray.map(({ label, desc, icon }: any, index: number) => <div
                                key={index}
                                className="flex al-stretch"
                            >
                                <div className="count-img">{icon}</div>
                                <div className="flex flex-col">
                                    <span className="fw-bold">{label}</span>
                                    <span className="primary-gray mn-fs">{desc}</span>
                                </div>
                            </div>)
                        }
                    </div>

                    <div className="bolck-content">
                        <Row gutter={22}>
                            <Col span={12}>
                                <div className="block-item fw-bold">
                                    <div className="block-title nm-fs">{firstOne?.name}</div>

                                    <div className="flex mgB16 mgT12">
                                        {
                                            firstOne?.sub_spaces?.map(({ name, id }: any) => (<div
                                                key={id}
                                                onClick={() => { setPar({ ...par, [firstOne.id]: id }) }}
                                                className={`block-txt mgR24 pointer ${par[firstOne.id] === id ? 'block-txt-active' : ''}`}>{name}</div>))
                                        }
                                    </div>
                                    {
                                        firstOne?.sub_spaces?.at(0)?.documents?.map((v: any) => (<div
                                            key={v}
                                            className="flex al-center space-between mgB24">
                                            <div className="news-title elli">集团举办“党建引领帮扶工作”资源大讲堂</div>
                                            <div className="news-txt white-nowrap flex primary-gray">
                                                <div className="mgR24">张某</div>
                                                <div className="mgL24">2025-09</div>
                                            </div>
                                        </div>))
                                    }
                                </div>
                            </Col>

                            <Col span={12}>
                                <div style={{ height: 'calc(100% - 24px)' }} className="flex flex-col">
                                    <Row gutter={22}>
                                        <Col span={12}>
                                            <div
                                                className="opera-box flex al-center space-between">
                                                <div
                                                    className="flex al-center">
                                                    <img src={IconFileAdd} alt="" />
                                                    <span className="nm-fs fw-bold">新增文档知识</span>
                                                </div>
                                                <img
                                                    onClick={() => { navigate('/knowledge/add') }}
                                                    className="pointer" src={IconRight} alt="" />
                                            </div>
                                        </Col>

                                        <Col span={12} className="flex flex-col">
                                            <div className="opera-box flex al-center space-between">
                                                <div className="flex al-center">
                                                    <img src={IconFloderAdd} alt="" />
                                                    <span className="nm-fs fw-bold">批量导入知识</span>
                                                </div>
                                                <img className="pointer" src={IconRight} alt="" />
                                            </div>
                                        </Col>
                                    </Row>

                                    <div className="flex1 label-cloud mgT24">
                                        <div style={{ color: '#010101' }} className="mgB16 nm-fs fw-bold">标签云</div>
                                        <div className="tag-box flex flex-wrap">
                                            {
                                                ['合同', '招标方案', '建筑法', '合同', '招标方案', '建筑法', '合同', '招标方案', '建筑法',]
                                                    .map((v: any, i: number) => (<div
                                                        className="tag-item mgR24 mn-fs pointer mgB12" key={i}>
                                                        {v}
                                                    </div>))
                                            }
                                        </div>
                                    </div>
                                </div>

                            </Col>

                            {
                                lastItems.map(({ name, id, sub_spaces }: any) => (<Col key={id} span={12}>
                                    <div className="block-item fw-bold">
                                        <div className="block-title nm-fs">{name}</div>

                                        <div className="flex mgB16 mgT12">
                                            {sub_spaces.map(({ id: subId, name: subName }: any) => <div
                                                key={subId}
                                                onClick={() => { setPar({ ...par, [id]: subId }) }}
                                                className={`block-txt mgR24 pointer ${par[id] === subId ? 'block-txt-active' : ''}`}>{subName}</div>)}
                                        </div>
                                        {
                                            getDocument(id)?.map((v: any) => (<div
                                                key={v}
                                                className="flex al-center space-between mgB24">
                                                <div className="news-title elli">集团举办“党建引领帮扶工作”资源大讲堂</div>
                                                <div className="news-txt white-nowrap flex primary-gray">
                                                    <div className="mgR24">张某</div>
                                                    <div className="mgL24">2025-09</div>
                                                </div>
                                            </div>))
                                        }
                                    </div>
                                </Col>))
                            }
                        </Row>
                    </div>
                </React.Fragment>
            }

            {
                documents.length !== 0 && <div className="documents-box">
                    {documents.map((v: any) => (<div
                        key={v.chunk_id}
                        className="document-item pointer">
                        <div className="hg-fs fw-bold mgB6">{v.title}</div>
                        <div>{v.content}</div>
                    </div>))}
                </div>
            }

            {isEmpty && <Empty style={{ marginTop: '10%' }} />}

        </div>
    )
}
