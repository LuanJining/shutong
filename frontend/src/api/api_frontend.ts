import { Par_Change_Pwd, Par_Check_Permission, Par_Classes, Par_Common_Params, Par_Search, Par_Space, Par_Sub_Space, Par_Upload_File, Par_Users } from "@/types/api";
import { del, get, post, put } from "@/utils/http";

/** @returns 登录*/
const login = (par: { login: string; password: string; }): Promise<any> => post('/auth/login', par);

/** @returns 创建用户*/
const createUser = (par: Par_Users): Promise<any> => post('/users', par);

/** @returns 用户列表*/
const getUsers = (par: Par_Common_Params): Promise<any> => get('/users', par);

/** @returns 分配角色*/
const assignRoles = (par: { roles: string[], userId: string ,space_id:number}): Promise<any> => post(`/spaces/${par.space_id}/members`, {
    userId: par.userId, roles: par.roles
});

/** @returns 获取用户角色*/
const getUserRoles = (userId: string): Promise<any> => get(`/spaces/members/${userId}`);

/** @returns 角色列表 */
const getRoles = (par: Par_Common_Params): Promise<any> => get('/roles', par);

/** @return 获取某个角色的权限*/
const getRolePermissions = (roleId: string): Promise<any> => get(`/roles/${roleId}/permissions`);

/** @returns 权限列表 */
const getPermissions = (): Promise<any> => get('/permissions');

/** @returns 创建的空间信息 */
const createSpace = (par: Par_Space): Promise<any> => post('/spaces', par);

/** @returns 空间列表 */
const getSpaces = (par?: Par_Common_Params): Promise<any> => get('/spaces', par);

/** @return 获取空间详情 */
const getSpaceById = (spaceId: string): Promise<any> => get(`/spaces/${spaceId}`);

/** @returns 权限检查结果 */
const checkPermission = (params: Par_Check_Permission): Promise<any> => post('/permissions/check', params);

/** @return 更新空间 */
const updateSpace = (spaceId: string, par: Par_Space): Promise<any> => put(`/spaces/${spaceId}`, par);

/** @return 删除空间 */
const deleteSpace = (spaceId: string): Promise<any> => del(`/spaces/${spaceId}`);

/** @return 获取单个用户*/
const getUserById = (userId: string): Promise<any> => get(`/users/${userId}`);

/** @return 获取单个角色 */
const getRoleById = (roleId: string): Promise<any> => get(`/roles/${roleId}`);

/** @return 获取单个权限 */
const getPermissionById = (permissionId: string): Promise<any> => get(`/permissions/${permissionId}`);

/** @return 提交审批流程 */
const uploadFile = (par: Par_Upload_File): Promise<any> => post(`/documents/upload`, par,
    {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
);

/** @return 获取上传的文档流 */
const getFile = (documentId: string | number): Promise<any> => get(`/documents/${documentId}/preview`, {}, {}, {
    responseType: 'blob'
});

/** @return 获取待处理任务 */
const getTasks = (par: Par_Common_Params): Promise<any> => get(`/workflow/tasks`, par);

/** @return 审批任务 */const taskOpear = (par: {
    task_id: number, comment: string, status: string
}): Promise<any> => post(`/workflow/tasks/approve`, par);

/** @return 用户看审批进度 */
const userTasks = (): Promise<any> => get(`/workflow/instances/user`);

/** @return 文档详情 */
const documentDetail = (documentId: string | number): Promise<any> => get(`/documents/${documentId}/info`);

/** @return 文档列表 */
const documentList = (space_id: string | number, par: Par_Common_Params): Promise<any> => get(`/documents/space/${space_id}`, par);

/** @return 修改密码 */
const changePwd = (par: Par_Change_Pwd): Promise<any> => post(`/auth/change-password`, par);

/** @return 主页搜索 */
const search = (par: Par_Search): Promise<any> => post(`/documents/search`, par);

/** @return 创建二级知识空间 */
const addSubSpaces = (par: Par_Sub_Space): Promise<any> => post(`/spaces/sub-spaces`, par);

/** @return 创建知识分类 */
const addClasses = (par: Par_Classes): Promise<any> => post(`/spaces/classes`, par);

/** @return 首页数据 */
const homePage = (): Promise<any> => get(`/documents/homepage`);

/** @return 启用 / 不启用 */
const docOpera = ({ documentId, status }: { documentId: number, status: 'publish' | 'unpublish' }): Promise<any> => post(`/documents/${documentId}/${status}`);

/** @return 获取标签云 */
const getTags = (): Promise<any> => get(`/documents/tag-cloud`);

export default {
    login, createUser, assignRoles, getUsers, getRoles, getPermissions, createSpace, getSpaces,
    getSpaceById, checkPermission, updateSpace, deleteSpace, getRolePermissions, getUserById, getRoleById,
    getPermissionById, uploadFile, getFile, getTasks, taskOpear, userTasks, documentDetail, documentList,
    changePwd, search, addSubSpaces, addClasses, homePage, docOpera, getTags,getUserRoles
}

