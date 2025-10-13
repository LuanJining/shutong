import _cache from "@/config/caches";
import storage from "@/utils/storage";
import { createSlice } from "@reduxjs/toolkit";

const initialState: {
    userInfo: any,
    isLogin: boolean,
    isLoading: boolean,
} = {
    isLogin: storage.get(_cache.AUTH_INFO)?.access_token,
    userInfo: storage.get(_cache.AUTH_INFO)?.user,
    isLoading: false,
};

export const systemSlice = createSlice({
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
        setIsLoading(state, { payload }) {
            state.isLoading = payload.isLoading;
        },
    },
});

export const { setUserInfo, setIsLogin, setIsLoading } = systemSlice.actions;

export default systemSlice.reducer; 
