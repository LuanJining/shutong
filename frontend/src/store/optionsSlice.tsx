import api_frontend from "@/api/api_frontend";
import _cache from "@/config/caches";
import { OPTIONS_TYPE } from "@/types/common";
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";

const initialState: {
    konwledges: OPTIONS_TYPE[],
} = {
    konwledges: []
};

export const initKonwledges = createAsyncThunk('optionsSlice/konwledges', async () => {
    const { data: { spaces } } = await api_frontend.getSpaces()
    return spaces.map(({ name, id }: any) => ({ label: name, value: id }))
})

export const optionsSlice = createSlice({
    //创建 slice
    name: "optionsSlice",
    initialState,
    reducers: {
        //可执行的reducer
        setKonwledges(state, { payload }) {
            state.konwledges = payload.konwledges;
        },
    },
    extraReducers(builder) {
        builder.addCase(initKonwledges.fulfilled, (state, action) => {
            state.konwledges = action.payload;
        });
    },
});

export const { setKonwledges } = optionsSlice.actions;

export default optionsSlice.reducer; 
