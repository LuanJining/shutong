import { get, post, put, del } from "@/utils/http";

/** @returns 登录*/
const login = (par: { login: string; password: string; }): Promise<any> => post('/iam/auth/login', par);

export default {
    login
}
