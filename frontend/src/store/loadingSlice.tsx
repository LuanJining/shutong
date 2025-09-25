import {  createSlice } from "@reduxjs/toolkit";

// export const loadingAsync = createAsyncThunk(
//   "loadingSlice/loadingThunk",
//   async (payload: any, thunkApi) => {
//     thunkApi.dispatch(setIsLoading({ isLoading: false }));
//   }
// );

export const loadingSlice = createSlice({
  //创建 slice
  name: "loadingSlice",
  initialState: {
    //状态初始值
    isLoading: false,
  },
  reducers: {
    //可执行的reducer
    setIsLoading(state, { payload }) {
      state.isLoading = payload.isLoading;
    },
  },
//   extraReducers: (builder) => {
//     builder
//       .addCase(loadingAsync.pending, (state, action) => {})
//       .addCase(loadingAsync.fulfilled, (state, action) => {})
//       .addCase(loadingAsync.rejected, (state, action) => {});
//   },
});

export const { setIsLoading } = loadingSlice.actions; //暴露方法，方便在需要使用的页面引入

export default loadingSlice.reducer; //暴露出reducer，创建store时需要
