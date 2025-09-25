// 系统配置初始化
// export const getSystemInfo = createAsyncThunk("getSystemInfo", async (_payload, thunkApi) => {
//     const r: any[] = await Promise.all([
//         store.dispatch(updateLabels()),
//         store.dispatch(initIndustrys()),
//     ])
//     const isReady: boolean = !r.some((v: any) => v.meta.requestStatus === 'rejected')
//     thunkApi.dispatch(setIsReady({ isReady: isReady ? 'ready' : 'error' }))
// })

import { createSlice } from "@reduxjs/toolkit";

const initialState: {
    isLogin: boolean,
    isReady: 'prepare' | 'ready' | 'error',
} = {
    isReady: 'prepare',
    isLogin: false
};

export const systemSlice = createSlice({
    //创建 slice
    name: "systemSlice",
    initialState,
    reducers: {
        //可执行的reducer
        // setIsReady(state, { payload }) {
        //     state.isReady = payload.isReady;
        // },
    },
});

export const {  } = systemSlice.actions;

export default systemSlice.reducer; 
