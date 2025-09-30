import utils from ".";
import storage from "./storage";
import _caches from "@/config/_caches";
import axios, { AxiosResponse, AxiosRequestConfig } from "axios";
import { message } from "antd";
import { getViteUrl } from "./tools";

const instance = axios.create({
    baseURL: getViteUrl('VITE_API_URL'),
    timeout: 120000,
});

instance.interceptors.request.use(
    (config) => {
        const token = storage.get(_caches.AUTH_INFO)?.access_token
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
            !config.url?.includes('/iam') && (config.headers['X-User-ID'] = storage.get(_caches.AUTH_INFO)?.user?.id)
        }
        return config;
    },
    (err) => {
        return Promise.reject(err);
    }
);

instance.interceptors.response.use(
    (res) => {
        return res;
    },
    (err) => {
        console.log(err)
        message.error(err?.response?.data?.error ?? err?.message)
        utils.setLoading(false)
        return Promise.reject(err);
    },
);
// 定义一个通用的请求函数
async function request<T>(
    config: AxiosRequestConfig
): Promise<AxiosResponse<T>> {
    try {
        const response = await instance.request<T>(config);
        return response;
    } catch (error) {
        throw error;
    }
}
// 封装GET请求
export async function get<T>(
    url: string,
    params?: Record<string, any>,
    headers?: any,
    config?: any
): Promise<T> {
    try {
        const response = await request<T>({
            method: "get",
            url,
            params,
            headers,
            ...config
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}
// 封装POST请求
export async function post<T>(
    url: string,
    data?: any,
    headers?: any,
    config?: any
): Promise<T> {
    try {
        const response = await request<T>({
            method: "post",
            url,
            data,
            headers,
            ...config
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}
// 封装PUT请求
export async function put<T>(
    url: string,
    data?: any,
    headers?: any
): Promise<T> {
    try {
        const response = await request<T>({
            method: "put",
            url,
            data,
            headers,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}
// 封装DELETE请求
export async function del<T>(
    url: string,
    params?: Record<string, any>,
    headers?: any
): Promise<T> {
    try {
        const response = await request<T>({
            method: "delete",
            url,
            params,
            headers,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

// 封装PATCH请求
export async function pat<T>(
    url: string,
    data?: any,
    headers?: any
): Promise<T> {
    try {
        const response = await request<T>({
            method: "patch",
            url,
            data,
            headers,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}