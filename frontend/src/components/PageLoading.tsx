import "./styles/page-loading.scss";
import React from "react";
import { Spin } from "antd";
import { useSelector } from "react-redux";

const PageLoading: React.FC = () => {
    const loading = useSelector((state: any) => state.systemSlice.isLoading)
    if (!loading) return null

    return (
        <div className="loading-all flex-center">
            <Spin size="large" />
        </div>
    );
};

export default PageLoading;
