import dayjs from 'dayjs'
import IconUser from '@/assets/icons/icon-user.png'
import { useEffect, useMemo, useState } from 'react'
import { useSelector } from "react-redux"
import api_frontend from '@/api/api_frontend'

export default function Personal() {
    const [permis, setPers] = useState<string>('')
    const userInfo: any = useSelector((state: any) => state.systemSlice.userInfo) ?? {}

    useEffect(() => { userInfo?.roles?.[0]?.id && getPermissions() }, [userInfo])

    const rolesInfo = useMemo(() => Object.entries(userInfo?.roles ?? []).map(([_key, obj]: any) => {
        return `${obj.display_name}（${obj.description}）`
    }), [userInfo])

    const getPermissions = async () => {
        const { permissions }: any = await api_frontend.getRolePermissions(userInfo?.roles?.[0]?.id)
        setPers(permissions.map((v: any) => v.display_name).join('、'))
    }


    return (
        <div className="app-common-deal">
            <div className="personal-item">
                <div className="personal-title">用户头像：</div>
                <div className="personal-desc">
                    <img style={{ width: 40, height: 40, objectFit: 'cover', borderRadius: '50%' }} src={userInfo.avatar ? userInfo.avatar : IconUser} alt="" />
                </div>
            </div>

            <div className="personal-item">
                <div className="personal-title">用户名称：</div>
                <div className="personal-desc">{userInfo.username}</div>
            </div>

            <div className="personal-item">
                <div className="personal-title">所属公司：</div>
                <div className="personal-desc">{userInfo.company}</div>
            </div>

            <div className="personal-item">
                <div className="personal-title">所属部门：</div>
                <div className="personal-desc">{userInfo.department}</div>
            </div>

            <div className="personal-item">
                <div className="personal-title">电子邮箱：</div>
                <div className="personal-desc">{userInfo.email}</div>
            </div>
            <div className="personal-item">
                <div className="personal-title">联系电话：</div>
                <div className="personal-desc">{userInfo.phone}</div>
            </div>

            <div className="personal-item" style={{ alignItems: 'flex-start' }}>
                <div className="personal-title">用户角色：</div>
                <div className="personal-desc">
                    {
                        rolesInfo.map((v: string, i: number) => (<div className='mgB12' key={i}>{v}</div>))
                    }
                </div>
            </div>


            <div className="personal-item" style={{ alignItems: 'flex-start' }}>
                <div className="personal-title">拥有权限：</div>
                <div className="personal-desc">{permis} </div>
            </div>

            <div className="personal-item">
                <div className="personal-title">最后登录时间：</div>
                <div className="personal-desc">{dayjs(userInfo.last_login).format('YYYY-MM-DD HH:mm')}</div>
            </div>

        </div>
    )
}
