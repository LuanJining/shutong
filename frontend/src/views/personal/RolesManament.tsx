import api_frontend from "@/api/api_frontend";
import _options from "@/config/_options";
import { Par_Common_Params } from "@/types/api";
import { Col, Input, Row, Space, Table, TableProps } from "antd";
import { useEffect, useState } from "react";

interface DataType {
    id: string;
    name: string;
    display_name: number;
    description: string;
    permissions: any[]
}

export default function RolesManament() {
    const [data, setData] = useState<DataType[]>([])
    const [par, setPar] = useState<Par_Common_Params>({ page: 1, page_size: 10, total: 0 })

    useEffect(() => { getRoles() }, [])

    const getRoles = async (values: any = {}) => {
        const params: any = { ...par, ...values }
        delete params?.total
        const { data: { roles, pagination: { total } } }: any = await api_frontend.getRoles(params)
        setData(roles.map((v: any) => ({ key: v.id, ...v })))
        setPar({ ...params, total })
    }

    const columns: TableProps<DataType>['columns'] = [
        {
            title: '序列号',
            dataIndex: 'id',
        },
        {
            title: '角色名称',
            dataIndex: 'name',
        },
        {
            title: '角色别称',
            dataIndex: 'display_name',
        },
        {
            title: '角色描述',
            dataIndex: 'description',
        },
        {
            title: '权限',
            dataIndex: 'permissions',
            width: 500,
            render: (permissions: any[]) => permissions.map((v: any) => v.display_name).join('、')
        },

        {
            title: '操作',
            width: 200,
            render: (item: any) => (
                <Space size="middle">
                    {/* <span>同意</span> */}
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

            <Table<DataType>
                columns={columns}
                dataSource={data}
                className="mgT24"
                pagination={{
                    current: +par.page,
                    pageSize: +par.page_size,
                    total: +(par.total ?? 0),
                    onChange: (page: number) => getRoles({ page })
                }}
            />
        </div>
    )
}
