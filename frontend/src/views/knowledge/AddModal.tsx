import utils from "@/utils";
import _opts from "@/config/opts"
import api_frontend from "@/api/api_frontend";
import { Form, Input, message, Modal } from "antd";
import { useForm } from "antd/es/form/Form";
import { Par_Space } from "@/types/api";
import { useEffect, useMemo } from "react";
import _ from "lodash";

interface IProps {
    open: boolean, setOpen: Function, callback: Function, item: any
}

export default function AddModal({ open, setOpen, callback, item }: IProps) {
    const [form] = useForm()

    useEffect(() => {
        item?.type === 'edit' && open && form.setFieldsValue(item)
        !open && form.resetFields()
    }, [item, open])

    const isNew: boolean = useMemo(() => item?.type?.includes('add'), [item])

    const isSubspace: boolean = useMemo(() => item?.type?.includes('subspace'), [item])

    const onFinish = async (values: Par_Space) => {

        try {
            utils.setLoading(true)

            if (item?.type === 'subspace-add') {
                await api_frontend.addSubSpaces({ ...values, space_id: +item?.space_id })
            }
            if (item?.type === 'classes-add') {
                await api_frontend.addClasses({ ...values, sub_space_id: +item?.space_id })
            }
            message.success(`${isNew ? '新增' : '编辑'}成功`)
            utils.setLoading(false)
            setOpen(false)
            callback()
        } catch (e) {
            utils.setLoading(false)
            throw (e)
        }
    }

    return (
        <Modal
            title={`${isNew ? '新增' : '编辑'}${isSubspace ? '子空间' : '知识分类'}`}
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
                labelCol={{ span: 5 }}
                wrapperCol={{ offset: 1 }}
                autoComplete="off"
                className="mgT24"
            >
                <Form.Item
                    label={`${isSubspace ? '子空间名称' : '知识分类'}`}
                    name="name"
                    rules={[{ required: true, message: `请输入${isSubspace ? '子空间名称' : '知识分类'}` }]}
                >
                    <Input placeholder={`请输入${isSubspace ? '子空间名称' : '知识分类'}`} />
                </Form.Item>

                <Form.Item
                    label={`${isSubspace ? '子空间描述' : '知识分类描述'}`}
                    name="description"
                    rules={[{ required: true, message: `请输入${isSubspace ? '子空间描述' : '知识分类描述'}` }]}
                >
                    <Input.TextArea
                        rows={4}
                        style={{ resize: 'none' }}
                        placeholder={`请输入${isSubspace ? '子空间描述' : '知识分类描述'}`} />
                </Form.Item>
            </Form>
        </Modal>
    )
}
