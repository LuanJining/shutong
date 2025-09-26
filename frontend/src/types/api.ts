export interface Par_Users {
    username: string;
    phone: string;
    email: string;
    password: string;
    nickname: string;
    department: string;
    company: string;
}

export interface Par_Space {
    name: string;
    description: string;
    type: string;
}

export interface Par_Check_Permission {
    space_id: number | string;
    resource: string;
    action: string;
}