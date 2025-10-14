import { useEffect, useRef } from "react";

const PagePdf = ({ page, scale = 1.2 }: { page: any; scale?: number }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!canvasRef.current || !page) return;

    const canvas = canvasRef.current;
    const context = canvas.getContext('2d');
    if (!context) return;

    const viewport = page.getViewport({ scale });

    canvas.width = viewport.width;
    canvas.height = viewport.height;

    page.render({
      canvasContext: context,
      viewport,
    }).promise.then(() => {
    }).catch((err:any) => {
      console.error('Page render error:', err);
    });

  }, [page, scale]);

  return <canvas ref={canvasRef} />
};

export default PagePdf