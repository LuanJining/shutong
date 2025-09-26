import _cache from "@/config/_caches";
import storage from "@/utils/storage";
import { createSlice } from "@reduxjs/toolkit";

const initialState: {
    isLogin: boolean,
    userInfo: any,
} = {
    isLogin: storage.get(_cache.AUTH_INFO)?.access_token,
    userInfo: storage.get(_cache.AUTH_INFO)?.user,
};

export const systemSlice = createSlice({
    //创建 slice
    name: "systemSlice",
    initialState,
    reducers: {
        //可执行的reducer
        setUserInfo(state, { payload }) {
            state.userInfo = payload.userInfo;
        },
        setIsLogin(state, { payload }) {
            state.isLogin = payload.isLogin;
        },
    },
});

export const { setUserInfo, setIsLogin } = systemSlice.actions;

export default systemSlice.reducer; 
