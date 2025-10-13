import "./styles/cate-select.scss"
import api_frontend from "@/api/api_frontend";
import { message, Modal } from "antd";
import { RightOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";

interface IProps {
    open: boolean;
    setOpen: Function;
    callback: (par: any) => void
}

export default function CateSelect({ open, setOpen, callback }: IProps) {
    const [spaces, setSpaces] = useState<any[]>([])
    const [par, setPar] = useState<any>({})

    useEffect(() => {
        open && getSpaces()
    }, [open])

    const getSpaces = async () => {
        const { data: { spaces } }: any = await api_frontend.getSpaces()
        setSpaces(spaces)
    }

    const sub_spaces: any[] = useMemo(() => spaces.find(({ id }: any) => id === par?.space?.id)?.sub_spaces ?? [], [spaces, par?.space])
    const classes: any[] = useMemo(() => sub_spaces.find(({ id }: any) => id === par?.subSpace?.id)?.classes ?? [], [sub_spaces, par?.subSpace])

    return (
        <Modal
            open={open}
            centered
            title="选择类别"
            className="cate-modal"
            width='50%'
            onCancel={() => { setOpen(false) }}
            onOk={() => {
                if (!par?.space?.id || !par?.subSpace?.id || !par?.classes?.id) {
                    message.error('请完善分类选择')
                    return
                }
                callback(par)
                setOpen(false)
            }}
        >
            <div className="cate-content flex mgT24">
                <div className="cate-box flex space-between">
                    <div className="cate-inner">
                        {
                            spaces.map(({ id, name }: any) => (<div
                                key={id}
                                onClick={() => { setPar({ space: { id, name } }) }}
                                className={`cate-item ${par?.space?.id === id ? 'cate-active' : ''}`}>{name}</div>))
                        }
                    </div>
                    <div className="flex-center"><RightOutlined className="lg-fs pointer" /></div>
                </div>
                <div className="cate-box flex space-between">
                    <div className="cate-inner">
                        {
                            sub_spaces?.map(({ id, name }: any) => (<div
                                key={id}
                                onClick={() => { setPar({ ...par, subSpace: { id, name }, classes: {} }) }}
                                className={`cate-item ${par?.subSpace?.id === id ? 'cate-active' : ''}`}>{name}</div>))
                        }

                    </div>
                    <div className="flex-center"><RightOutlined className="lg-fs pointer" /></div>
                </div>
                <div className="cate-box">
                    <div className="cate-inner">
                        {
                            classes?.map(({ id, name }: any) => (<div
                                key={id}
                                onClick={() => { setPar({ ...par, classes: { id, name } }) }}
                                className={`cate-item ${par?.classes?.id === id ? 'cate-active' : ''}`}>{name}</div>))
                        }
                    </div>
                </div>
            </div>
        </Modal>
    )
}
