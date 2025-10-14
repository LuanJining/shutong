import dayjs from 'dayjs'
import { Button, Col, Input, Row, Space, Table, TableProps } from "antd";
import CreatUser from "./CreatUser";
import { useEffect, useState } from "react";
import api_frontend from "@/api/api_frontend";
import AssignRoles from './AssignRoles';
import { Par_Common_Params } from '@/types/api';

interface DataType {
    username: string;
    phone: string;
    email: number;
    company: string;
    nickname: string;
    department: string;
    created_at: string;
}

export default function UsersManagement() {
    const [open, setOpen] = useState<boolean>(false)
    const [assign, setAssign] = useState<string>('')
    const [data, setData] = useState<DataType[]>([])
    const [par, setPar] = useState<Par_Common_Params>({ page: 1, page_size: 10, total: 0 })

    useEffect(() => { getUsers() }, [])

    const getUsers = async (values: any = {}) => {
        const params: any = { ...par, ...values }
        delete params?.total
        const { data: { users, pagination: { total } } } = await api_frontend.getUsers(params)
        setData(users.map((v: any) => ({ key: v.id, ...v })))
        setPar({ ...params, total })
    }


    const columns: TableProps<DataType>['columns'] = [
        {
            title: '用户名',
            dataIndex: 'username',
        },

        {
            title: '昵称',
            dataIndex: 'nickname',
        },

        {
            title: '所属公司 / 部门',
            render: (item: any) => `${item.company} / ${item.department}`
        },
        {
            title: '邮箱',
            dataIndex: 'email',
        },

        {
            title: '电话',
            dataIndex: 'phone',
        },

        {
            title: '创建时间',
            dataIndex: 'created_at',
            render: (created_at: string) => dayjs(created_at).format('YYYY-MM-DD HH:mm')
        },
        {
            title: '操作',
            render: (item: any) => (
                <Space size="middle">
                    <span onClick={() => { setAssign(item.id) }} className='pointer primary-blue'>分配角色</span>
                    {/* <span>驳回</span> */}
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
                <Col className="flex jf-end" offset={16} span={3}>
                    <Button onClick={() => { setOpen(true) }} type="primary">新增用户</Button>
                </Col>
            </Row>

            <Table<DataType>
                columns={columns}
                dataSource={data}
                className="mgT24"
                pagination={{
                    current: +par.page,
                    pageSize: +par.page_size,
                    total: +(par.total ?? 0),
                    onChange: (page: number) => getUsers({ page })
                }}
            />

            <CreatUser open={open} setOpen={setOpen} callback={() => { getUsers() }} />
            <AssignRoles open={Boolean(assign)} setOpen={setAssign} callback={() => { }} userId={assign} />
        </div>
    )
}
