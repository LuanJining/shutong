import "./index.scss"
import dayjs from 'dayjs'
import api_frontend from '@/api/api_frontend'
import { ArrowLeftOutlined } from "@ant-design/icons"
import { Button, Col, Input, message, Popconfirm, Row, Space, Table, TableProps } from "antd"
import { useEffect, useState } from 'react'
import { useNavigate } from "react-router-dom"
import _optsEnum from "@/config/_optsEnum"
import AddModal from "./AddModal"
import utils from "@/utils"

interface DataType {
    id: string;
    name: string;
    description: string;
    created_at: string;
    type: string
}

export default function Index() {
    const navigate = useNavigate()
    const [list, setList] = useState<any[]>([])
    const [open, setOpen] = useState<boolean>(false)
    const [curItem, setCurItem] = useState<any>(null)

    useEffect(() => { getSpaces() }, [])

    const getSpaces = async () => {
        try {
            const { data: { spaces } }: any = await api_frontend.getSpaces()
            setList(spaces.map((v: any) => ({ key: v.id, ...v })))
        } catch (e) {
            throw (e)
        }
    }

    const onDelete = async (item: any) => {
        try {
            utils.setLoading(true)
            await api_frontend.deleteSpace(item.id)
            await getSpaces()
            message.success('删除成功')
            utils.setLoading(false)
        } catch (e) {
            utils.setLoading(false)
            throw (e)
        }

    }

    console.log(list)

    const columns: TableProps<DataType>['columns'] = [
        {
            title: '序列号',
            dataIndex: 'id',
        },
        {
            title: '知识库名称',
            dataIndex: 'name',
        },
        {
            title: '知识库类型',
            dataIndex: 'type',
            render: (type: keyof typeof _optsEnum.SPACE_TYPE) => _optsEnum.SPACE_TYPE[type]
        },
        {
            title: '知识库描述',
            dataIndex: 'description',
            width: 400,
            ellipsis: true
        },
        {
            title: '创建人',
            dataIndex: 'creator',
            render: (creator: any) => creator.username
        },
        {
            title: '创建时间',
            dataIndex: 'created_at',
            render: (created_at: string) => dayjs(created_at).format('YYYY-MM-DD HH:mm')
        },


        {
            title: '操作',
            width: 200,
            render: (item: any) => (
                <Space size="middle">
                    <span onClick={() => { setOpen(true); setCurItem(item) }} className="pointer primary-blue">编辑</span>
                    <Popconfirm
                        title="删除知识库"
                        description="确认要删除知识库吗？"
                        okText="确认"
                        cancelText="取消"
                        onConfirm={() => { onDelete(item) }}
                    >
                        <span className="pointer primary-blue">删除</span>
                    </Popconfirm>

                </Space>
            ),
        },
    ];

    return (
        <div className='app-knowledge'>
            <div onClick={() => { navigate(-1) }} className="pointer">
                <ArrowLeftOutlined className="ls-fs mgR6" />
                <span>返回</span>
            </div>

            <div className="knowledge-box">
                <div className="nm-fs fw-bold">知识库管理</div>
                <Row className="mgT24">
                    <Col span={4}><Input placeholder="请输入关键字" /></Col>

                    <Col className="text-right" offset={16} span={4}>
                        <Button
                            onClick={() => {
                                setOpen(true)
                                setCurItem(null)
                            }}
                            type="primary">新增知识库</Button></Col>
                </Row>

                <Table<DataType>
                    columns={columns}
                    dataSource={list}
                    className="mgT24"
                />
            </div>

            <AddModal open={open} setOpen={setOpen} callback={() => { getSpaces() }} item={curItem} />
        </div>
    )
}
