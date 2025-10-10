import dayjs from 'dayjs'
import utils from '@/utils';
import _optsEnum from '@/config/_optsEnum';
import api_frontend from '@/api/api_frontend';
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom';
import { Par_Common_Params } from '@/types/api';
import { DeleteOutlined, DownloadOutlined } from '@ant-design/icons';
import { Button, Col, message, Popconfirm, Row, Space, Switch, Table, TableProps } from 'antd';

interface DataType {
    id: string;
    name: string;
    description: string;
    created_at: string;
    type: string
}

export default function DocumentManagement({ space_id }: { space_id: string | number }) {
    const navigate = useNavigate()
    const [list, setList] = useState<any[]>([])
    const [par, setPar] = useState<Par_Common_Params>({ page: 1, page_size: 10 })

    useEffect(() => { space_id && getDocument() }, [space_id])

    const getDocument = async (values: any = {}) => {
        try {
            const params: Par_Common_Params = { ...par, ...values }
            delete params?.total
            const { data: { items, total } }: any = await api_frontend.documentList(space_id, params)
            setList(items.map((v: any) => ({ key: v.id, ...v })))
            setPar({ ...params, total })
        } catch (e) {
            throw (e)
        }
    }

    const download = async (documentId: string, filename: string) => {
        utils.setLoading(true)
        const res: any = await api_frontend.getFile(documentId)
        utils.downloadFromFlow(res, filename)
        utils.setLoading(false)
    }

    const onDelete = async (item: any) => {
        try {
            utils.setLoading(true)
            await api_frontend.deleteSpace(item.id)
            await getDocument()
            message.success('删除成功')
            utils.setLoading(false)
        } catch (e) {
            utils.setLoading(false)
            throw (e)
        }
    }

    const columns: TableProps<DataType>['columns'] = [
        {
            title: '名称',
            dataIndex: 'title',
            width: '60%',
            ellipsis: true,
            className: 'pointer',
            render: (title: string, data: any) => <div
                onClick={() => { navigate("/document/detail", { state: { documentId: data?.id } }) }}
            >{title}</div>
        },
        {
            title: '上传人',
            dataIndex: 'created_by',
        },
        {
            title: '上传时间',
            dataIndex: 'created_at',
            render: (created_at: string) => dayjs(created_at).format('YYYY-MM-DD HH:mm')
        },

        {
            title: '属性',
            dataIndex: 'file_type',
        },
        {
            title: '操作',
            width: 150,
            align: 'center',
            render: (item: any) => (
                <Space size='large' className='flex al-center jf-center'>
                    <Switch size='small' defaultChecked />
                    <DownloadOutlined onClick={() => { download(item.id, item.file_name) }} className='pointer lg-fs' style={{ color: '#4190FF' }} />
                    <Popconfirm
                        title="删除文档"
                        description="确认要删除当前文档吗？"
                        okText="确认"
                        cancelText="取消"
                        onConfirm={() => { onDelete(item) }}
                    >
                        <DeleteOutlined className='pointer lg-fs' style={{ color: '#BD3124' }} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (

        <div className="knowledge-box pd16">
            <div className="nm-fs fw-bold">文档</div>
            <Row className="mgT24">
                <div className="sm-fs primary-gray">知识库的所有文件都在这里显示，整个知识库都可以链接到应用引用或通过 Chat 插件进行索引。</div>
                <Col className="text-right" flex={1}>
                    <Button
                        style={{ width: 150 }}
                        onClick={() => { navigate('/knowledge/add') }}
                        type="primary">
                        + 添加文档
                    </Button></Col>
            </Row>

            <Table<DataType>
                columns={columns}
                dataSource={list}
                className="mgT24"
                scroll={{ x: 'max-content' }}
                pagination={{
                    current: +par.page,
                    pageSize: +par.page_size,
                    total: +(par.total ?? 0),
                    onChange: (page: number) => getDocument({ page })
                }}
            />
        </div>
    )
}
