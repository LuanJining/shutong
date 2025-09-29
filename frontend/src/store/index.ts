import { configureStore } from "@reduxjs/toolkit";

import optionsSlice from "./optionsSlice";
import systemSlice from "./systemSlice";

export default configureStore({
    reducer: {
        systemSlice, optionsSlice
    },
});
