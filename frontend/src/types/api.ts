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
    visibility: 'public' | 'internal' | 'private' | 'protected';
    urgency?: 'normal' | 'urgent';
    tags: string;
    summary: string;
    created_by: string;
    department: string;
    need_approval: boolean
}

export interface Par_Chat {
    space_id: string;
    question:string;
    document_ids: [];
    limit: number;
}

export interface Par_Change_Pwd {
    old_password: string;
    new_password:string;
}