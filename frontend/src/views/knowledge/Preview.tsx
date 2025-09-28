import _ from 'lodash';
import mammoth from 'mammoth';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import { useState, useEffect, useMemo } from 'react';
import { usePDFStreamRenderer } from '@/hooks/usePDFStreamRenderer';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl;

export default function CustomFileViewer({ fileInfo }: any) {
    const fileType: string = fileInfo?.fileType
    const [docxHtml, setDocxHtml] = useState<string>('');
    const source: any = useMemo(() => {
        if (!fileInfo?.file || fileType !== 'pdf') return null;
        return { type: 'file', file: fileInfo.file };
    }, [fileInfo?.file])
    const { pdfPages } = usePDFStreamRenderer(source)

    useEffect(() => {
        !_.isEmpty(fileInfo) && handleFileChange()
    }, [fileInfo])

    const handleFileChange = () => {
        const selectedFile = fileInfo.file
        setDocxHtml('');

        if (!selectedFile) return;

        if (fileType === 'docx') {
            renderDocx(selectedFile);
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
        <div className='pdf-container'>
            {fileType === 'pdf' && pdfPages.length !== 0
                ? pdfPages
                : fileType === 'docx' && docxHtml ? (
                    <div className='pdL24 pdR24 pdT16 pdT16' dangerouslySetInnerHTML={{ __html: docxHtml }}
                    />
                ) : <></>}
        </div>
    );
};