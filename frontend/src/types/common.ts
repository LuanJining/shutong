import React from "react"

export interface OPTIONS_TYPE {
    label: any,
    value: any
}

export interface ENUM_TYPE {
    [key: string]: any
}

export interface ROUTESELF_TYPE {
    path?: string;
    element?: React.ReactNode;
    meta?: {
        title?: string;
        icon?: React.ReactNode;
        show?: boolean; // 是否展示在nav中
    };
    children?: ROUTESELF_TYPE[];
}
