import React, { useState, useRef, useEffect } from 'react';
import mammoth from 'mammoth';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import _ from 'lodash';

pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl;

export default function CustomFileViewer({ fileInfo }: any) {
    const [pdfPages, setPdfPages] = useState<JSX.Element[] | null>(null);
    const [docxHtml, setDocxHtml] = useState<string>('');

    const fileType: string = fileInfo?.fileType

    useEffect(() => {
        !_.isEmpty(fileInfo) && handleFileChange()
    }, [fileInfo])

    console.log(fileType)

    const handleFileChange = () => {
        const selectedFile = fileInfo.file
        setPdfPages(null);
        setDocxHtml('');

        if (!selectedFile) return;


        if (!['pdf', 'docx'].includes(fileType || '')) {
            console.log('只支持 PDF 和 DOCX 文件');
            return;
        }

        if (fileType === 'pdf') {
            render(selectedFile);
        } else if (fileType === 'docx') {
            renderDocx(selectedFile);
        }
    };
    // ========== PDF 渲染：使用 PDF.js ==========
    const render = async (file: any) => {
        try {
            const arrayBuffer = await file.arrayBuffer();
            const pdf = await pdfjsLib.getDocument(arrayBuffer).promise;

            const renderedPages: JSX.Element[] = [];

            for (let i = 1; i <= pdf.numPages; i++) {
                const page = await pdf.getPage(i);
                const viewport = page.getViewport({ scale: 2 });

                const canvas = document.createElement('canvas');
                const context = canvas.getContext('2d');
                if (!context) continue;

                canvas.width = viewport.width;
                canvas.height = viewport.height;

                await page.render({
                    canvasContext: context,
                    viewport,
                }).promise;

                renderedPages.push(
                    <div key={i}>
                        <canvas  ref={(el) => {
                            if (el) el.replaceWith(canvas);
                            // 或者更简单的：直接 appendChild(canvas) 到某个 DOM 容器中
                        }} />
                    </div>
                );
            }

            setPdfPages(renderedPages);
        } catch (err) {
            console.error('PDF 渲染失败:', err);
        }
    };

    // ========== DOCX 渲染：使用 mammoth.js ==========
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
            {fileType === 'pdf' && pdfPages ? (
                <div>{pdfPages}</div>
            ) : fileType === 'docx' && docxHtml ? (
                <div className='pdL24 pdR24 pdT16 pdT16' dangerouslySetInnerHTML={{ __html: docxHtml }}
                />
            ) : <></>}
        </div>
    );
};