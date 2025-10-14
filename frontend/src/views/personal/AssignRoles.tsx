import "./assign-modal.scss"
import _ from "lodash";
import utils from "@/utils";
import _opts from "@/config/opts"
import api_frontend from "@/api/api_frontend";
import { useEffect, useState } from "react";
import { Checkbox, message, Modal } from "antd";
import { OPTIONS_TYPE } from './../../types/common';

interface IProps {
    open: boolean, setOpen: Function, callback: Function, userId: string
}

interface Props_Choosed {
    [key: string]: string[]
}

export default function AssignRoles({ open, setOpen, callback, userId }: IProps) {
    const [spaces, setSpaces] = useState<any[]>([])
    const [chooseInfo, setCInfo] = useState<Props_Choosed>({})
    const [activeSpace, setActiveSpace] = useState<number>(-1)

    useEffect(() => { open && getSpaces(); userId && getUserRoles() }, [open])

    const getUserRoles = async () => {
        const { data }: any = await api_frontend.getUserRoles(userId)
        const tempCInfo: any = {}
        data?.map(({ roles, space_id }: any) => { tempCInfo[space_id] = roles })
        setCInfo(tempCInfo)
    }

    const getSpaces = async () => {
        const { data: { spaces: interSpaces } }: any = await api_frontend.getSpaces({ page: 1, page_size: 50 })
        setSpaces(interSpaces)
        setActiveSpace(interSpaces?.[0]?.id ?? -1)
    }

    const onFinish = async () => {
        if (chooseInfo?.[activeSpace]?.length === 0) {
            message.success('请先选中需要分配的角色')
            return
        }
        utils.setLoading(true)
        await api_frontend.assignRoles({ userId, roles: chooseInfo?.[activeSpace], space_id: activeSpace })
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
            className="assign-modal"
            width={640}
        >
            <div className="flex assign-box">
                <div className="space-box pd24 mgR24">
                    {spaces.map(({ id, name }: any) => (<div
                        key={id}
                        className={`space-item ${activeSpace === id ? 'active-space-item' : ''}`}
                    >{name}</div>))}
                </div>
                <div className="roles-box pd24" style={{ paddingLeft: 48 }}>
                    {
                        _opts.ROLES.map(({ label, value }: OPTIONS_TYPE) => {
                            return (<div
                                key={value}
                                onClick={(e) => {
                                    e.stopPropagation()
                                    e.preventDefault()
                                    const curCArray: any[] = chooseInfo?.[activeSpace] ?? []
                                    console.log(curCArray)
                                    setCInfo({
                                        ...chooseInfo,
                                        [activeSpace]: curCArray.includes(value)
                                            ? _.difference(curCArray, [value])
                                            : [...curCArray, value]
                                    })
                                }}
                                className="role-item mgB12">
                                <Checkbox checked={chooseInfo?.[activeSpace]?.includes(value)}>
                                    <span className="mgL12">{label}</span>
                                </Checkbox>
                            </div>)
                        })
                    }
                </div>
            </div>
        </Modal>
    )
}
