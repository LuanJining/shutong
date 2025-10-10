import "./index.scss"
import { useEffect, useState } from 'react'
import _optsEnum from "@/config/_optsEnum"
import _opts from '@/config/_opts';
import api_frontend from "@/api/api_frontend"
import DocumentManagement from "./DocumentManagement"
import PageNav from "@/components/PageNav"
import { Props_Self_Nav } from "@/types/pages";

export default function Index() {
    const [pathKey, setPathKey] = useState<string | number>('')
    const [items, setItems] = useState<(Props_Self_Nav & any)[]>([])

    useEffect(() => { getSpaces() }, [])
    useEffect(() => { setPathKey(items?.[0]?.key ?? '') }, [items])

    const getSpaces = async () => {
        const { data: { spaces } }: any = await api_frontend.getSpaces()
        setItems(spaces.map((v: any) => ({ ...v, label: v.name, key: v.id })))
    }

    return (
        <div className='app-knowledge'>
            <PageNav pathKey={pathKey} pathArray={items} setPathKey={setPathKey} getSpaces={getSpaces}/>
            <div className="nav-content flex1">
                <DocumentManagement space_id={pathKey} />
            </div>
        </div>
    )
}
