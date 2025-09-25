import { configureStore } from "@reduxjs/toolkit";

import loadingSlice from "./loadingSlice";
import systemSlice from "./systemSlice";

export default configureStore({
  reducer: {
    loadingSlice,systemSlice
  },
});
