import { get, post, put, del } from "@/utils/http";

/** @returns 登录*/
const login = (par: { username: string; password: string; }): Promise<any> => post('/users/login', par);

export default {
    login
}
