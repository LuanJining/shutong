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
import { Col, Input, Row } from "antd"

const countArray: any[] = [
    { label: '制度汇总', desc: '公司各部门制定指引', icon: <img src={IconFile} />, },
    { label: '法律法规', desc: '法律法规资源库', icon: <img src={IconResume} />, },
    { label: '项目经验', desc: '企业项目案例', icon: <img src={IconTarget} />, },
    { label: '行政后勤', desc: '后勤职能部门资源库', icon: <img src={IconFloder} />, },
    { label: '产品研发', desc: '企业产品档案库', icon: <img src={IconMouse} />, },
]

export default function Index() {
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
                    <Input variant="borderless" placeholder="请输入搜索内容" />
                    <SearchOutlined className="hg-fs" />
                </div>
            </div>


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
                            <div className="block-title nm-fs">集团文档案</div>

                            <div className="flex mgB16 mgT12">
                                <div className="block-txt mgR24">公文</div>
                                <div className="block-txt mgR24">新闻</div>
                                <div className="block-txt">党政学习</div>
                            </div>

                            {
                                [1, 2, 3, 4, 5, 6].map((v: any) => (<div
                                    key={v}
                                    className="flex al-center space-between mgB24">
                                    <div className="news-title elli">集团举办“党建引领帮扶工作”资源大讲堂</div>
                                    <div className="news-txt flex primary-gray">
                                        <div className="mgR24">张某</div>
                                        <div className="mgL24">2025-09</div>
                                    </div>
                                </div>))
                            }
                        </div>
                    </Col>

                    <Col span={12}>
                        <div style={{height:'calc(100% - 24px)'}} className="flex flex-col">
                            <Row gutter={22}>
                                <Col span={12}>
                                    <div className="opera-box flex al-center space-between">
                                        <div className="flex al-center">
                                            <img src={IconFileAdd} alt="" />
                                            <span className="nm-fs fw-bold">新增文档知识</span>
                                        </div>
                                        <img className="pointer" src={IconRight} alt="" />
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
                                <div style={{color:'#010101'}} className="mgB16 nm-fs fw-bold">标签云</div>
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
                        [1, 2, 3, 4].map((v: any) => (<Col key={v} span={12}>
                            <div className="block-item fw-bold">
                                <div className="block-title nm-fs">集团文档案</div>

                                <div className="flex mgB16 mgT12">
                                    <div className="block-txt mgR24">公文</div>
                                    <div className="block-txt mgR24">新闻</div>
                                    <div className="block-txt">党政学习</div>
                                </div>

                                {
                                    [1, 2, 3, 4, 5, 6].map((v: any) => (<div
                                        key={v}
                                        className="flex al-center space-between mgB24">
                                        <div className="news-title elli">集团举办“党建引领帮扶工作”资源大讲堂</div>
                                        <div className="news-txt flex primary-gray">
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


        </div>
    )
}
