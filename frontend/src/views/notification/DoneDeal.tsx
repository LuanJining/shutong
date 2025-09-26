import { Col, Input, Row, Space, Table, TableProps } from "antd";

interface DataType {
    key: string;
    name: string;
    age: number;
    address: string;
    tags: string[];
}

export default function DoneDeal() {

    const data: DataType[] = [
        {
            key: '1',
            name: 'John Brown',
            age: 32,
            address: 'New York No. 1 Lake Park',
            tags: ['nice', 'developer'],
        },
        {
            key: '2',
            name: 'Jim Green',
            age: 42,
            address: 'London No. 1 Lake Park',
            tags: ['loser'],
        },
        {
            key: '3',
            name: 'Joe Black',
            age: 32,
            address: 'Sydney No. 1 Lake Park',
            tags: ['cool', 'teacher'],
        },
    ];

    const columns: TableProps<DataType>['columns'] = [
        {
            title: '标题',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '分类',
            dataIndex: 'age',
            key: 'age',
        },
        {
            title: '发起人',
            dataIndex: 'address',
            key: 'address',
        },

        {
            title: '接收时间',
            dataIndex: 'address',
            key: 'address',
        },
        {
            title: '处理时间',
            dataIndex: 'address',
            key: 'address',
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
            />
        </div>
    )
}
