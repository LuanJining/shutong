import api_frontend from '@/api/api_frontend';
import _optsEnum from '@/config/optsEnum';
import { Par_Common_Params } from '@/types/api';
import utils from '@/utils';
import { Button, Col, Input, message, Popconfirm, Row, Space, Table, TableProps } from 'antd';
import dayjs from 'dayjs';
import { useEffect, useState } from 'react';
import AddModal from './AddModal';

interface DataType {
    id: string;
    name: string;
    description: string;
    created_at: string;
    type: string
}

export default function KonwledgeManagement() {
    const [list, setList] = useState<any[]>([])
    const [open, setOpen] = useState<boolean>(false)
    const [curItem, setCurItem] = useState<any>(null)
    const [par, setPar] = useState<Par_Common_Params>({ page: 1, page_size: 10, total: 0 })

    useEffect(() => { getSpaces() }, [])

    const getSpaces = async (values: any = {}) => {
        try {
            const params: any = { ...par, ...values }
            delete params?.total
            const { data } = await api_frontend.getSpaces(params)
            setList(data.map((v: any) => ({ key: v.id, ...v })))
            setPar({ ...params, total: data.length })
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

    const columns: TableProps<DataType>['columns'] = [
        {
            title: '知识库名称',
            dataIndex: 'name',
        },
        {
            title: '知识库描述',
            dataIndex: 'description',
            width: '50%',
            ellipsis: true
        },
        {
            title: '知识库类型',
            dataIndex: 'type',
            render: (type: keyof typeof _optsEnum.SPACE_TYPE) => _optsEnum.SPACE_TYPE[type]
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
            width: 150,
            align: 'center',
            render: (item: any) => (
                <Space size="middle" className='flex-center'>
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

        <div className="knowledge-box pd16">
            <div className="nm-fs fw-bold">知识库</div>
            <Row className="mgT24">
                {/* <div className="sm-fs primary-gray">知识库的所有文件都在这里显示，整个知识库都可以链接到应用引用或通过 Chat 插件进行索引。</div> */}
                <Col span={4} ><Input placeholder='输入关键字' /></Col>
                <Col className="text-right" flex={1}>
                    <Button
                        style={{ width: 150 }}
                        onClick={() => {
                            setOpen(true)
                            setCurItem(null)
                        }}
                        type="primary">
                        + 新增知识库
                    </Button></Col>
            </Row>

            <Table<DataType>
                columns={columns}
                dataSource={list}
                className="mgT24"
                pagination={{
                    current: +par.page,
                    pageSize: +par.page_size,
                    total: +(par.total ?? 0),
                    onChange: (page: number) => getSpaces({ page })
                }}
            />
            <AddModal open={open} setOpen={setOpen} callback={() => { getSpaces() }} item={curItem} />
        </div>
    )
}
