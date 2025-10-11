export interface Props_File_View {
    file: File | null;
    type: 'file' | 'url';
    styles: {},
    fileType: 'pdf' | 'docx' | '';
}

export interface Props_Self_Nav {
    key: number | string; label: string
}