import base64
import io
import logging
import os
import threading
from typing import Dict, List, Optional

import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from paddleocr import PaddleOCR
from PIL import Image, UnidentifiedImageError

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
            _ocr_models[lang] = PaddleOCR(use_angle_cls=True, lang=lang, show_log=False)
    return _ocr_models[lang]


def _decode_image(data_base64: str) -> np.ndarray:
    try:
        decoded = base64.b64decode(data_base64)
    except base64.binascii.Error as exc:
        raise HTTPException(status_code=400, detail=f"Invalid base64 payload: {exc}") from exc

    try:
        image = Image.open(io.BytesIO(decoded))
        image = image.convert("RGB")
    except (UnidentifiedImageError, ValueError) as exc:
        raise HTTPException(status_code=400, detail=f"Unsupported image format: {exc}") from exc

    return np.array(image)


@app.get("/healthz")
def healthz() -> Dict[str, str]:
    return {"status": "ok"}


@app.post("/v1/ocr", response_model=OCRResponse)
def run_ocr(request: OCRRequest) -> OCRResponse:
    if not request.content_base64:
        raise HTTPException(status_code=400, detail="content_base64 is required")

    image = _decode_image(request.content_base64)
    ocr = _load_ocr(request.language or _default_language)

    results = ocr.ocr(image, det=True, rec=True)
    lines: List[str] = []
    for page in results:
        for item in page:
            if len(item) < 2:
                continue
            text = item[1][0]
            if text:
                lines.append(text)

    text_output = "\n".join(lines)
    if not text_output:
        raise HTTPException(status_code=422, detail="OCR did not return any text")

    return OCRResponse(text=text_output, lines=lines)
