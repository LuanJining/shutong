import utils from "@/utils";
import _opts from "@/config/opts"
import api_frontend from "@/api/api_frontend";
import { Form, Input, message, Modal, Select } from "antd";
import { useForm } from "antd/es/form/Form";
import { Par_Space } from "@/types/api";
import { useEffect, useMemo } from "react";
import _ from "lodash";

interface IProps {
    open: boolean, setOpen: Function, item: any;
    getSpaces: () => void
}

export default function AssignRoles({ open, setOpen, item, getSpaces }: IProps) {
    const [form] = useForm()

    useEffect(() => {
        item && open && form.setFieldsValue(item)
        !open && form.resetFields()
    }, [item, open])

    const isNew: boolean = useMemo(() => !Boolean(item), [item])

    const onFinish = async (values: Par_Space) => {

        try {
            utils.setLoading(true)
            await (isNew
                ? api_frontend.createSpace(values)
                : api_frontend.updateSpace(item.id, values))
            message.success(`${isNew ? '创建' : '编辑'}成功`)
            utils.setLoading(false)
            setOpen(false)
            getSpaces()
        } catch (e) {
            utils.setLoading(false)
            throw (e)
        }
    }

    return (
        <Modal
            title={`${isNew ? '创建' : '编辑'}知识库`}
            centered
            open={open}
            destroyOnClose
            onCancel={() => { setOpen('') }}
            onOk={() => { form.submit() }}
            width={580}
        >
            <Form
                form={form}
                onFinish={onFinish}
                labelCol={{ span: 4 }}
                wrapperCol={{ offset: 1 }}
                autoComplete="off"
                className="mgT24"
            >
                <Form.Item
                    label="知识库名称"
                    name="name"
                    rules={[{ required: true, message: '请输入知识库名称!' }]}
                >
                    <Input placeholder="请输入知识库名称" />
                </Form.Item>

                <Form.Item
                    label="知识库描述"
                    name="description"
                    rules={[{ required: true, message: '请输入知识库描述!' }]}
                >
                    <Input.TextArea
                        rows={4}
                        style={{ resize: 'none' }}
                        placeholder="请输入知识库描述" />
                </Form.Item>

                <Form.Item
                    label="类型"
                    name="type"
                    rules={[
                        { required: true, message: '请选择知识库类型' },
                    ]}
                >
                    <Select options={_opts.SPACE_TYPE} placeholder="请选择知识库类型" />
                </Form.Item>
            </Form>
        </Modal>
    )
}
