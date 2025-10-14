import "./styles/preview.scss"
import _ from 'lodash';
import mammoth from 'mammoth';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import { useState, useEffect, useMemo, useRef } from 'react';
import { usePDFStreamRenderer } from '@/hooks/usePDFStreamRenderer';
import { Props_File_View } from '@/types/pages';
import { useLocation } from 'react-router-dom';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl;

export default function CustomFileViewer({ fileType, file, type, styles }: Props_File_View) {
    const pdfRef = useRef<any>(null)
    const [scrollTop, setScrollTop] = useState<number>(0)
    const [docxHtml, setDocxHtml] = useState<string>('');

    const documentId: any = useLocation().state?.documentId

    const source: any = useMemo(() => {
        if (fileType !== 'pdf') return null;
        const result: any = { type, scrollTop }
        type === 'file'
            ? result.file = file
            : result.documentId = documentId
        return result
    }, [type, fileType, scrollTop])

    const { pdfPages, loading,totalPage } = usePDFStreamRenderer(source)

    console.log(pdfPages)

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

    const onScroll = (e: any) => {
        const sTop: number = e.target.scrollTop
        setScrollTop(sTop ?? 0)
    }

    return (
        <div ref={pdfRef}
            onScroll={onScroll}
            className='pdf-container' style={styles}>
            {fileType === 'pdf' && pdfPages.length !== 0
                ? <div style={{height:`${totalPage * 1100}px`}}>{pdfPages}</div>
                : fileType === 'docx' && docxHtml ? (
                    <div className='pdL24 pdR24 pdT16 pdT16' dangerouslySetInnerHTML={{ __html: docxHtml }}
                    />
                ) : <></>}
        </div>
    );
};