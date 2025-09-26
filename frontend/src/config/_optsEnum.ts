import { ENUM_TYPE } from "@/types/common";

const ROLES: ENUM_TYPE = {
    1: '超级管理员',// super_admin
    2: '企业管理员',// corp_admin
    3: '空间管理员',// space_admin
    4: '内容审核员',// content_reviewer
    5: '内容编辑者',// content_editor
    6: '只读用户'// read_only_user
}

const SPACE_TYPE: ENUM_TYPE = {
    department: '部门',
    project: '项目',
    team: '团队',
}

const optsEnum = {
    ROLES,SPACE_TYPE
};

type EnumOType = typeof optsEnum;

type DynamicType<T> = {
    [K in keyof T]: {
        [P in keyof T[K]]: T[K][P]
    }
};

type EnumODynamicType = DynamicType<EnumOType>;

export default optsEnum as EnumODynamicType;