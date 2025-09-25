import optsEnum from "./optsEnum";

// 定义 enumS 的类型
type EnumS = typeof optsEnum;

// 动态生成 temp 的类型
type TempType = {
    [K in keyof EnumS]: Array<{ label: EnumS[K][keyof EnumS[K]], value: keyof EnumS[K] }>;
};

// 初始化 temp，并指定动态类型
const temp: TempType = {} as TempType;

// 填充 temp
(Object.keys(optsEnum) as (keyof EnumS)[]).forEach((v: keyof EnumS) => {
    temp[v] = (Object.keys(optsEnum[v]) as any).map((v1: keyof EnumS[typeof v]) => ({
        label: optsEnum[v][v1],
        value: v1
    }));
});

export default temp as TempType;