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

/**
 * binary 文件封装
 * @param key:string
 * @param value:File
 * @returns formData
 */
export function getFormData(values: {}) {
    if (!values || _.isEmpty(values)) return values

    const formData: any = new FormData();
    Object.entries(values).map(([key, value]: any) => {
        formData.append(key, value);
    })
    return formData;
}

export default {
    setLoading, isFalseValue, normFile, getFormData
}


