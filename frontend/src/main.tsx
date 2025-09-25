
import 'dayjs/locale/zh-cn';
import "@/assets/styles/common.scss";
import "@/assets/styles/global.scss";
import "@/assets/fonts/iconfont/iconfont.css"

import dayjs from 'dayjs';
import App from "./App.tsx";
import store from "@/store";
import zhCN from 'antd/locale/zh_CN';
import ReactDOM from "react-dom/client";

import { ConfigProvider } from "antd";
import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";

dayjs.locale('zh-cn');

ReactDOM.createRoot(document.getElementById("root")!).render(
    <ConfigProvider locale={zhCN}>
        <Provider store={store}>
            <BrowserRouter>
                <App />
            </BrowserRouter>
        </Provider>
    </ConfigProvider>
);
