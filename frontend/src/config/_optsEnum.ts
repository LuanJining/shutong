import { ENUM_TYPE } from "@/types/common";

const TEST: ENUM_TYPE = {
    test: 'tests'
}

const optsEnum = {
    TEST
};

type EnumOType = typeof optsEnum;

type DynamicType<T> = {
    [K in keyof T]: {
        [P in keyof T[K]]: T[K][P]
    }
};

type EnumODynamicType = DynamicType<EnumOType>;

export default optsEnum as EnumODynamicType;