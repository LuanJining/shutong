export interface Par_Common_Params {
    page: string | number;
    page_size: string | number;
    total?: string | number;
}

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

export interface Par_Upload_File {
    file: File;
    file_name: string;
    space_id: number | string;
    sub_space_id: string;
    class_id: string;
    visibility: 'public' | 'internal' | 'private' | 'protected';
    urgency?: 'normal' | 'urgent';
    tags: string;
    summary: string;
    created_by: string;
    department: string;
    need_approval: boolean;
    version: string;
    use_type: string
}

export interface Par_Chat {
    space_id: string;
    question: string;
    document_ids: [];
    limit: number;
}

export interface Par_Change_Pwd {
    old_password: string;
    new_password: string;
}

export interface Par_Search {
    query?: string,
    space_id?: string | number,
    sub_space_id?: string | number,
    class_id?: string | number,
    limit?: string | number
}

export interface Par_Sub_Space {
    name: string;
    description: string;
    space_id: number;
}

export interface Par_Classes {
    name: string;
    description: string;
    sub_space_id: number;
}