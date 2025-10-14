import api_frontend from '@/api/api_frontend';
import PagePdf from '@/components/PagePdf';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import { useState, useEffect, useCallback } from 'react';

// 设置 PDF.js worker（必须在模块顶层调用！）
pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl

type RenderSource =
    | { type: 'file'; file: File; scrollTop: number }
    | { type: 'url'; documentId: number; scrollTop: number };

interface UsePDFStreamRendererReturn {
    pdfPages: JSX.Element[];
    loading: boolean;
    error: Error | null;
    totalPage:number
}

const CanvasWidth: number = 1100

export const usePDFStreamRenderer = (
    source: RenderSource | null = null,
): UsePDFStreamRendererReturn => {
    const [totalPage, setPages] = useState<number>(0)
    const [pdfObj, setPdf] = useState<any>(null)
    const [pdfPages, setPdfPages] = useState<JSX.Element[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => { source?.type && renderPDF() }, [source?.type]);

    const renderPDF = useCallback(async () => {
        if (!source) return;

        setLoading(true);
        setError(null);
        setPdfPages([]);

        try {
            let arrayBuffer: ArrayBuffer;

            if (source.type === 'file') {
                // 情况 1：传入的是 File 对象（比如用户上传）
                const file = source.file;
                arrayBuffer = await file.arrayBuffer();
            } else if (source.type === 'url') {
                const response = await api_frontend.getFile(source.documentId); // 假设这个方法等价于你下面的 getFile
                arrayBuffer = await response.arrayBuffer();
            } else {
                throw new Error('无效的 source 类型');
            }
            const pdf = await pdfjsLib.getDocument(arrayBuffer).promise;
            setPdf(pdf)
            setPages(pdf.numPages)
        } catch (err) {
            console.error('PDF 渲染失败:', err);
            setError(err instanceof Error ? err : new Error('未知错误'));
        } finally {
            setLoading(false);
        }
    }, [source?.type])

    useEffect(() => { getPages() }, [source?.scrollTop, totalPage])

    const getPages = async () => {
        const scrollTop = source?.scrollTop ?? 0
        const currentPage: number = Math.max(Math.ceil(scrollTop / CanvasWidth), 1)
        const end: number = Math.min(currentPage + 2, totalPage)
        const pdf = pdfObj

        if (!pdf) return

        const renderedPages: JSX.Element[] = [];
        for (let i = 1; i <= end; i++) {
            const page = await pdf.getPage(i);
            renderedPages.push(
                <div key={i} className='flex-center'>
                    <PagePdf page={page} />
                </div>
            );
        }
        setPdfPages(renderedPages);
    }

    return { pdfPages, loading, error ,totalPage};
};