


// 添加数据到缓存
export const addToCache = (key: string, value: any) => {
    localStorage.setItem(key, JSON.stringify(value));
};

// 从缓存中获取数据
export const getFromCache = (key: string) => {
    const value = localStorage.getItem(key);
    return value ? JSON.parse(value) : null;
};

// 从缓存中删除数据
export const removeFromCache = (key: string) => {
    localStorage.removeItem(key);
};

// 获取所有缓存数据
export const getAllFromCache = (prefix: string): Record<string, any> => {
    const keys = Object.keys(localStorage).filter(key => key.startsWith(prefix));
    const result: Record<string, number> = {};

    keys.forEach(key => {
        const value = localStorage.getItem(key);
        if (value !== null) {
            result[key] = JSON.parse(value);
        }
    });

    return result;
};
// 清除所有缓存数据
export const cleanAllFromCache = (prefix: string) => {
    const keys = Object.keys(localStorage).filter(key => key.startsWith(prefix));
    keys.forEach(key => removeFromCache(key))
};
