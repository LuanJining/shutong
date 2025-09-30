import utils from "@/utils";
import _opts from "@/config/_opts"
import api_frontend from "@/api/api_frontend";
import { useEffect, useState } from "react";
import { Checkbox, message, Modal } from "antd";
import { OPTIONS_TYPE } from './../../types/common';
import { useSelector } from "react-redux";

interface IProps {
    open: boolean, setOpen: Function, callback: Function, userId: string
}

export default function AssignRoles({ open, setOpen, callback, userId }: IProps) {
    // const userInfo = useSelector((state: any) => state.systemSlice.userInfo)
    const [choosed, setChoosed] = useState<number>(-1)

    // useEffect(() => { open && userInfo?.roles?.[0]?.id && setChoosed(userInfo?.roles?.[0]?.id.toString()) }, [userInfo, open])

    const onFinish = async () => {
        if (choosed === -1) {
            message.success('请先选中需要分配的角色')
            return
        }
        utils.setLoading(true)
        await api_frontend.assignRoles({ userId, role_id: choosed })
        message.success('分配成功')
        utils.setLoading(false)
        setOpen('')
        callback()
    }

    return (
        <Modal
            title='配置角色'
            centered
            open={open}
            destroyOnClose
            onCancel={() => { setOpen('') }}
            onOk={() => { onFinish() }}
        >
            <div className="roles-box pd24" style={{ paddingLeft: 48 }}>
                {
                    _opts.ROLES.map(({ label, value }: OPTIONS_TYPE) => {
                        return (<div
                            key={value}
                            onClick={() => { setChoosed(value) }}
                            className="role-item mgB12">
                            <Checkbox checked={value === choosed}>
                                <span className="mgL12">{label}</span>
                            </Checkbox>
                        </div>)
                    })
                }
            </div>
        </Modal>
    )
}
