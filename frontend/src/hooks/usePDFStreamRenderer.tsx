import api_frontend from '@/api/api_frontend';
import * as pdfjsLib from 'pdfjs-dist';
import pdfjsWorkerUrl from 'pdfjs-dist/build/pdf.worker.js?url';
import { useState, useEffect, useCallback } from 'react';

// 设置 PDF.js worker（必须在模块顶层调用！）
pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerUrl

type RenderSource =
    | { type: 'file'; file: File }
    | { type: 'url'; documentId: string | number; };

interface UsePDFStreamRendererReturn {
    pdfPages: JSX.Element[];
    loading: boolean;
    error: Error | null;
}

/**
 * Hook：用于渲染 PDF（支持传入 File 或远程 URL 接口）
 * @param source 渲染源：可以是 File 对象，也可以是 URL + 请求配置
 * @returns { pdfPages, loading, error }
 */
export const usePDFStreamRenderer = (
    source: RenderSource | null = null,
): UsePDFStreamRendererReturn => {
    const [pdfPages, setPdfPages] = useState<JSX.Element[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => { renderPDF() }, [source]);

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
            try {
                const pdf = await pdfjsLib.getDocument(arrayBuffer).promise;

                const renderedPages: JSX.Element[] = [];

                for (let i = 1; i <= pdf.numPages; i++) {
                    const page = await pdf.getPage(i);
                    renderedPages.push(
                        <div key={i} className='flex-center' style={{ marginBottom: 20 }}>
                            <canvas
                                ref={(el) => {
                                    if (!el) return;

                                    const canvasEl = el as HTMLCanvasElement;
                                    const context = canvasEl.getContext('2d');
                                    if (!context) return;

                                    const viewport = page.getViewport({ scale: 1.5 });

                                    canvasEl.width = viewport.width;
                                    canvasEl.height = viewport.height;

                                    page.render({
                                        canvasContext: context,
                                        viewport,
                                    }).promise;
                                }}
                            />
                        </div>
                    );
                }

                setPdfPages(renderedPages);
            } catch (e) {
                console.error('PDF 渲染失败:', e);
            }

        } catch (err) {
            console.error('PDF 渲染失败:', err);
            setError(err instanceof Error ? err : new Error('未知错误'));
        } finally {
            setLoading(false);
        }
    }, [source])

    return { pdfPages, loading, error };
};