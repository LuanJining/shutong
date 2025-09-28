import { configureStore } from "@reduxjs/toolkit";

import systemSlice from "./systemSlice";

export default configureStore({
    reducer: {
        systemSlice
    },
});
