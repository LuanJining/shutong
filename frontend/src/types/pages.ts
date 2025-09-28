export interface Porps_File_View {
    file: File | null;
    type: 'file' | 'url';
    styles:{},
    fileType: 'pdf' | 'docx' | '';
}