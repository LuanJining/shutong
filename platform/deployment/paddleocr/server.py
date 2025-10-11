import base64
import io
import logging
import os
import threading
from pathlib import Path
from typing import Dict, List, Optional, Sequence

import numpy as np
from docx import Document as DocxDocument
from fastapi import FastAPI, HTTPException
from paddleocr import PaddleOCR
from pdf2image import convert_from_bytes
from PIL import Image, UnidentifiedImageError
from pydantic import BaseModel

logger = logging.getLogger(__name__)

app = FastAPI(title="PaddleOCR Service", version="1.0.0")

_ocr_models: Dict[str, PaddleOCR] = {}
_ocr_lock = threading.RLock()
_default_language = os.getenv("PADDLE_LANG", "ch")


class OCRRequest(BaseModel):
    file_name: str
    content_base64: str
    language: Optional[str] = None


class OCRResponse(BaseModel):
    text: str
    lines: List[str]


def _load_ocr(language: str) -> PaddleOCR:
    lang = language or _default_language
    with _ocr_lock:
        if lang not in _ocr_models:
            logger.info("Loading PaddleOCR model for language: %s", lang)
            _ocr_models[lang] = PaddleOCR(
                use_angle_cls=True,
                lang=lang,
                show_log=False,
                use_gpu=False,
                enable_mkldnn=False,  # Disable MKL-DNN to avoid SIGILL
                cpu_threads=2,  # Limit CPU threads for stability
            )
    return _ocr_models[lang]


def _decode_base64(data_base64: str) -> bytes:
    try:
        return base64.b64decode(data_base64)
    except base64.binascii.Error as exc:
        raise HTTPException(status_code=400, detail=f"Invalid base64 payload: {exc}") from exc


def _image_from_bytes(data: bytes) -> np.ndarray:
    try:
        image = Image.open(io.BytesIO(data))
        image = image.convert("RGB")
    except (UnidentifiedImageError, ValueError) as exc:
        raise HTTPException(status_code=400, detail=f"Unsupported image format: {exc}") from exc

    return np.array(image)


def _render_pdf_to_images(data: bytes) -> List[np.ndarray]:
    try:
        pil_images = convert_from_bytes(data)
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=400, detail=f"Failed to convert PDF to images: {exc}") from exc

    if not pil_images:
        raise HTTPException(status_code=422, detail="PDF had no renderable pages")

    return [np.array(img.convert("RGB")) for img in pil_images]


def _extract_docx_text(data: bytes) -> List[str]:
    try:
        document = DocxDocument(io.BytesIO(data))
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=400, detail=f"Failed to read DOCX file: {exc}") from exc

    lines: List[str] = []
    for paragraph in document.paragraphs:
        text = paragraph.text.strip()
        if text:
            lines.append(text)

    for table in document.tables:
        for row in table.rows:
            cells = [cell.text.strip() for cell in row.cells if cell.text.strip()]
            if cells:
                lines.append("\t".join(cells))

    if not lines:
        raise HTTPException(status_code=422, detail="DOCX document did not contain readable text")

    return lines


def _run_ocr_on_images(images: Sequence[np.ndarray], language: str) -> List[str]:
    if not images:
        raise HTTPException(status_code=422, detail="No pages or images to OCR")

    ocr = _load_ocr(language or _default_language)
    lines: List[str] = []
    for image in images:
        results = ocr.ocr(image, det=True, rec=True)
        for page in results:
            for item in page:
                if len(item) < 2:
                    continue
                text = item[1][0]
                if text:
                    lines.append(text)

    if not lines:
        raise HTTPException(status_code=422, detail="OCR did not return any text")

    return lines


@app.get("/healthz")
def healthz() -> Dict[str, str]:
    return {"status": "ok"}


@app.post("/v1/ocr", response_model=OCRResponse)
def run_ocr(request: OCRRequest) -> OCRResponse:
    if not request.content_base64:
        raise HTTPException(status_code=400, detail="content_base64 is required")

    file_bytes = _decode_base64(request.content_base64)
    suffix = Path(request.file_name).suffix.lower()
    language = request.language or _default_language

    if suffix in {".jpg", ".jpeg", ".png", ".bmp", ".tif", ".tiff"}:
        images = [_image_from_bytes(file_bytes)]
        lines = _run_ocr_on_images(images, language)
        return OCRResponse(text="\n".join(lines), lines=lines)

    if suffix == ".pdf":
        images = _render_pdf_to_images(file_bytes)
        lines = _run_ocr_on_images(images, language)
        return OCRResponse(text="\n".join(lines), lines=lines)

    if suffix == ".docx":
        lines = _extract_docx_text(file_bytes)
        return OCRResponse(text="\n".join(lines), lines=lines)

    raise HTTPException(status_code=400, detail=f"Unsupported file type: {suffix}")
