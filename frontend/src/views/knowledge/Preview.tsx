import "./styles/preview.scss"
import _ from 'lodash';
import mammoth from 'mammoth';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import { useState, useEffect, useMemo } from 'react';
import { usePDFStreamRenderer } from '@/hooks/usePDFStreamRenderer';
import { Props_File_View } from '@/types/pages';
import { useLocation } from 'react-router-dom';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl;

export default function CustomFileViewer({ fileType, file, type, styles }: Props_File_View) {
    const [docxHtml, setDocxHtml] = useState<string>('');

    const documentId: any = useLocation().state?.documentId

    const source: any = useMemo(() => {
        if (fileType !== 'pdf') return null;
        const result: any = { type }
        type === 'file'
            ? result.file = file
            : result.documentId = documentId
        return result
    }, [type, fileType])

    const { pdfPages } = usePDFStreamRenderer(source)

    useEffect(() => {
        fileType === 'docx' && handleFileChange()
    }, [fileType])

    const handleFileChange = () => {
        setDocxHtml('');
        if (!file) return;

        if (fileType === 'docx') {
            renderDocx(file);
        }
    };

    const renderDocx = (file: File) => {
        const reader = new FileReader();
        reader.onload = async (event) => {
            const arrayBuffer = event.target?.result as ArrayBuffer;

            try {
                const result = await mammoth.convertToHtml({ arrayBuffer });
                const html = result.value;
                setDocxHtml(html);
            } catch (err) {
                console.log('无法解析 DOCX 文件');
                console.error(err);
            }
        };
        reader.readAsArrayBuffer(file);
    };

    return (
        <div className='pdf-container' style={styles}>
            {fileType === 'pdf' && pdfPages.length !== 0
                ? pdfPages
                : fileType === 'docx' && docxHtml ? (
                    <div className='pdL24 pdR24 pdT16 pdT16' dangerouslySetInnerHTML={{ __html: docxHtml }}
                    />
                ) : <></>}
        </div>
    );
};