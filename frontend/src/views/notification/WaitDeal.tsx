import dayjs from 'dayjs'
import api_frontend from "@/api/api_frontend";
import IconAgree from "@/assets/icons/icon-agree.png"
import IconRefuse from "@/assets/icons/icon-refuse.png"
import { Par_Common_Params } from "@/types/api";
import { Col, Input, message, Popconfirm, Row, Space, Table, TableProps } from "antd";
import { useEffect, useState } from "react";
import utils from '@/utils';

export default function WaitDeal() {
    const [list, setList] = useState<any[]>([])
    const [par, setPar] = useState<Par_Common_Params>({ page: 1, page_size: 10, total: 0 })
    const [comment, setComment] = useState<string>('')
    useEffect(() => { getList() }, [])

    const getList = async (values: any = {}) => {
        const params: any = { ...par, ...values }
        delete params?.total
        const { data: { items } }: any = await api_frontend.getTasks(params)
        setList(items.map((v: any) => ({ key: v.id, ...v })))
    }

    const opera = async (type: 'agree' | 'refuse', data: any) => {
        if (type === 'refuse') {
            message.info('敬请期待')
            return
        }

        if (!comment.trim()) {
            message.warning('请输入审批意见！')
            return
        }

        utils.setLoading(true)
        await api_frontend.taskAgree(data.id, comment)
        await getList()
        message.success('操作成功')
        setComment('')
        utils.setLoading(false)
    }

    const columns: TableProps<any>['columns'] = [
        {
            title: '标题',
            render: (data: any) => data.instance.title
        },
        {
            title: '描述',
            render: (data: any) => data.instance.description,
            width: 400,
            ellipsis: true
        },

        {
            title: '发起人',
            render: (data: any) => data.instance.created_by
        },

        {
            title: '接收时间',
            dataIndex: ' assigned_at',
            render: (assigned_at: any) => dayjs(assigned_at).format('YYYY-MM-DD HH:mm')
        },
        {
            title: '操作',
            render: (data: any) => (
                <Space size="middle">
                    <Popconfirm
                        placement="leftBottom"
                        title={'同意审批'}
                        description={<div className='flex flex-col white-nowrap mgB12'>
                            <span className='mgB6'>评审意见：</span>
                            <Input.TextArea
                                value={comment}
                                onChange={(e: any) => { setComment(e.target.value) }}
                                style={{ resize: 'none' }} rows={2} />
                        </div>}
                        okText="确认"
                        cancelText="取消"
                        onCancel={() => { setComment('') }}
                        onConfirm={() => { opera('agree', data) }}
                    >
                        <div className='flex al-center pointer'>
                            <img className='mgR6' style={{ width: 20, height: 20, objectFit: "cover" }} src={IconAgree} alt="" />
                            <span>同意</span>
                        </div>
                    </Popconfirm>
                    <Popconfirm
                        placement="leftBottom"
                        title={'驳回审批'}
                        description={'确定驳回当前审批吗？'}
                        okText="确认"
                        cancelText="取消"
                        onConfirm={() => { opera('refuse', data) }}
                    >
                        <div className='flex al-center pointer'>
                            <img className='mgR6' style={{ width: 20, height: 20, objectFit: "cover" }} src={IconRefuse} alt="" />
                            <span>驳回</span>
                        </div>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <div className="app-common-deal">
            <Row>
                <Col span={5}>
                    <Input placeholder="请输入内容" />
                </Col>
            </Row>

            <Table
                columns={columns}
                dataSource={list}
                className="mgT24"
            />
        </div>
    )
}
