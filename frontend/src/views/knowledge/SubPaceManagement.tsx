import "./styles/app-subspaces.scss"
import AddModal from './AddModal';
import _optsEnum from '@/config/optsEnum';
import { Button, } from 'antd';
import { useEffect, useMemo, useState } from "react";
interface IProps {
    getSpaces: () => void;
    space_id: number;
    curSpace: any
}

export default function SubPaceManagement({ space_id, getSpaces, curSpace }: IProps) {
    const [curItem, setCur] = useState<any>(null)
    const [open, setOpen] = useState<boolean>(false)
    const [subSpaceId, setSubSpaceId] = useState<string>('')

    useEffect(() => { setSubSpaceId(curSpace?.sub_spaces?.at(0)?.id) }, [curSpace?.sub_spaces])

    const classes: any[] = useMemo(() => curSpace?.sub_spaces?.find(({ id }: any) => id === subSpaceId)?.classes, [curSpace, subSpaceId])
    return (
        <div className="app-subspaces flex">
            <div>
                <Button
                    className="mg16"
                    style={{ width: 150 }}
                    onClick={() => {
                        setCur({ type: 'subspace-add', space_id })
                        setOpen(true)
                    }}
                    type="primary">
                    + 新增子空间
                </Button>
                <div className='subspaces-box'>
                    {curSpace?.sub_spaces?.map(({ name, id }: any) => (<div
                        key={id}
                        onClick={() => { setSubSpaceId(id) }}
                        className={`subspace-item ${id === subSpaceId ? 'space-active' : ''}`}>
                        {name}</div>))}
                </div>
            </div>

            <div>
                <Button
                    className="mg16"
                    style={{ width: 150, visibility: subSpaceId ? 'visible' : 'hidden' }}
                    onClick={() => {
                        setOpen(true)
                        setCur({ type: 'classes-add', space_id: subSpaceId })
                    }}
                    type="primary">
                    + 新增知识分类
                </Button>
                <div className='subspaces-box'>
                    {classes?.map(({ name, id }: any) => (<div
                        key={id}
                        className='subspace-item'>
                        {name}</div>))}
                </div>
            </div>
            <AddModal open={open} setOpen={setOpen} callback={() => { getSpaces() }} item={curItem} />
        </div>
    )
}
