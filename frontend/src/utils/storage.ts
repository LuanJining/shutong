export default {
    set,
    get,
    remove,
    clear
}
// 存储数据到 sessionStorage  
function set(key: string, value: any) {
    try {
        sessionStorage.setItem(key, JSON.stringify(value));
    } catch (error) {
        console.error('存储数据到 sessionStorage 时出错：', error);
    }
}

// 从 sessionStorage 中获取数据  
function get(key: string) {
    try {
        const item: any = sessionStorage.getItem(key);
        return item ? JSON.parse(item) : null
    } catch (error) {
        console.error('从 sessionStorage 中获取数据时出错：', error);
        return null;
    }
}

// 从 sessionStorage 中删除数据  
function remove(key: string) {
    try {
        sessionStorage.removeItem(key);
    } catch (error) {
        console.error('从 sessionStorage 中删除数据时出错：', error);
    }
}

// 从 sessionStorage 中删除数据  
function clear() {
    try {
        sessionStorage.clear()
    } catch (error) {
        console.error('删除所有信息出错', error);
    }
}  