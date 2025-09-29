import { get, post, put, del } from "@/utils/http";
import { Par_Check_Permission, Par_Common_Params, Par_Space, Par_Upload_File, Par_Users } from "@/types/api";

/** @returns 登录*/
const login = (par: { login: string; password: string; }): Promise<any> => post('/iam/auth/login', par);

/** @returns 创建用户*/
const createUser = (par: Par_Users): Promise<any> => post('/iam/users', par);

/** @returns 用户列表*/
const getUsers = (par: Par_Common_Params): Promise<any> => get('/iam/users', par);

/** @returns 分配角色*/
const assignRoles = (par: { role_id: number, userId: string }): Promise<any> => post(`/iam/users/${par.userId}/roles`, { role_id: par.role_id });

/** @returns 角色列表 */
const getRoles = (par: Par_Common_Params): Promise<any> => get('/iam/roles', par);

/** @param 获取某个角色的权限*/
const getRolePermissions = (roleId: string): Promise<any> => get(`/iam/roles/${roleId}/permissions`);

/** @returns 权限列表 */
const getPermissions = (): Promise<any> => get('/iam/permissions');

/** @returns 创建的空间信息 */
const createSpace = (par: Par_Space): Promise<any> => post('/iam/spaces', par);

/** @returns 空间列表 */
const getSpaces = (par?: Par_Common_Params): Promise<any> => get('/iam/spaces', par);

/** @param 获取空间详情 */
const getSpaceById = (spaceId: string): Promise<any> => get(`/iam/spaces/${spaceId}`);

/** @returns 权限检查结果 */
const checkPermission = (params: Par_Check_Permission): Promise<any> => post('/iam/permissions/check', params);

/** @param 更新空间 */
const updateSpace = (spaceId: string, par: Par_Space): Promise<any> => put(`/iam/spaces/${spaceId}`, par);

/** @param 删除空间 */
const deleteSpace = (spaceId: string): Promise<any> => del(`/iam/spaces/${spaceId}`);

/** @param 获取单个用户*/
const getUserById = (userId: string): Promise<any> => get(`/iam/users/${userId}`);

/** @param 获取单个角色 */
const getRoleById = (roleId: string): Promise<any> => get(`/iam/roles/${roleId}`);

/** @param 获取单个权限 */
const getPermissionById = (permissionId: string): Promise<any> => get(`/iam/permissions/${permissionId}`);

/** @param 提交审批流程 */
const uploadFile = (par: Par_Upload_File): Promise<any> => post(`/kb/upload`, par,
    {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
);

/** @param 获取上传的文档流 */
const getFile = (documentId: string | number): Promise<any> => get(`/kb/${documentId}/preview`, {}, {}, {
    responseType: 'blob'
});

/** @param 获取待处理任务 */
const getTasks = (par: Par_Common_Params): Promise<any> => get(`/workflow/tasks`, par);

/** @param 审批任务 */
const taskAgree = (taskId: string | number, comment: string): Promise<any> => post(`/workflow/tasks/${taskId}/approve`, { comment });

/** @param 用户看审批进度 */
const userTasks = (): Promise<any> => get(`/workflow/instances/user`);

/** @param 文档详情 */
const documentDetail = (documentId: string | number): Promise<any> => get(`/kb/${documentId}/info`);

/** @param 文档列表 */
const documentList = (space_id: string | number, par: Par_Common_Params): Promise<any> => get(`/kb/${space_id}/space`, par);

export default {
    login, createUser, assignRoles, getUsers, getRoles, getPermissions, createSpace, getSpaces,
    getSpaceById, checkPermission, updateSpace, deleteSpace, getRolePermissions, getUserById, getRoleById,
    getPermissionById, uploadFile, getFile, getTasks, taskAgree, userTasks, documentDetail, documentList
}

