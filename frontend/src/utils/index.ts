import _ from "lodash";
import store from "@/store";
import { setIsLoading } from "@/store/systemSlice";

/**
 * @returns 是否为 {} [] '' NAN  null undefined
 */
function isFalseValue(value: any) {
    const type: string = Object.prototype.toString.call(value);
    if (type === "[object String]") {
        return value === "";
    }
    if (type === "[object Array]") {
        return !Boolean(value.length);
    }
    if (type === "[object Number]") {
        return isNaN(value);
    }
    if (type === "[object Object]") {
        return !Boolean(Object.keys(value).length);
    }
    if (type === "[object Null]") {
        return true;
    }
    if (type === "[object Undefined]") {
        return true;
    }
}

/**
 * 设置全局加载状态
 * @param isLoading 加载状态
 */
function setLoading(isLoading: boolean) {
    store.dispatch(setIsLoading({ isLoading }));
}

function normFile(e: any) {
    if (Array.isArray(e)) {
        return e;
    }
    return e?.fileList;
};

/** @returns formData*/
function getFormData(values: { [key: string]: any }) {
    if (!values || _.isEmpty(values)) return values

    const formData: any = new FormData();
    Object.entries(values).map(([key, value]: any) => {
        formData.append(key, value);
    })
    return formData;
}

function downloadFromFlow(res: any, filename: string) {
    let blob = new Blob([res], {
        type: "application/octet-stream;charset=UTF-8",
    });
    let downloadElement: any = document.createElement("a");
    downloadElement.download = filename; // 文件名称 自定义
    downloadElement.href = window.URL.createObjectURL(blob);
    downloadElement.click();
    document.body.appendChild(downloadElement);
    document.body.removeChild(downloadElement);
    window.URL.revokeObjectURL(downloadElement.href);
}


export default {
    setLoading, isFalseValue, normFile, getFormData, downloadFromFlow
}