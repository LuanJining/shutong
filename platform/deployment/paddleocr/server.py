# -*- coding: utf-8 -*-
# server.py — PaddleOCR 服务（兼容多版本）
# IMPORTANT: environment variables MUST be set before importing numpy/paddle/paddleocr

import os

# -----------------------
# Conservative environment settings — must be set before any heavy imports
# -----------------------
# Try to disable MKL / oneDNN aggressive optimizations early
os.environ['FLAGS_use_mkldnn'] = 'false'
os.environ['FLAGS_use_dnnl'] = 'false'
os.environ['FLAGS_eager_delete_tensor_gb'] = '0.0'

# Restrict BLAS / threading to single thread for maximum compatibility
os.environ['CPU_NUM_THREADS'] = '1'
os.environ['OMP_NUM_THREADS'] = '1'
os.environ['OPENBLAS_NUM_THREADS'] = '1'
os.environ['MKL_NUM_THREADS'] = '1'
os.environ['NUMEXPR_NUM_THREADS'] = '1'

# Allow duplicate lib loading if necessary (keep your original behavior)
os.environ.setdefault('KMP_DUPLICATE_LIB_OK', 'TRUE')
os.environ.setdefault('MKL_THREADING_LAYER', 'GNU')

# -----------------------
# Now safe to import other modules
# -----------------------
import base64
import io
import logging
import threading
from pathlib import Path
from typing import Dict, List, Optional, Sequence

import inspect
import numpy as np
from docx import Document as DocxDocument
from fastapi import FastAPI, HTTPException

# Deliberately import paddleocr lazily inside _load_ocr to reduce chances of early native init
from pdf2image import convert_from_bytes
from PIL import Image, UnidentifiedImageError
from pydantic import BaseModel

logger = logging.getLogger("server")
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s:%(name)s: %(message)s")

app = FastAPI(title="PaddleOCR Service", version="1.0.0")

# store loaded PaddleOCR instances per language
_ocr_models: Dict[str, object] = {}
_ocr_lock = threading.RLock()
_default_language = os.getenv("PADDLE_LANG", "ch")


class OCRRequest(BaseModel):
    file_name: str
    content_base64: str
    language: Optional[str] = None


class OCRResponse(BaseModel):
    text: str
    lines: List[str]


def _load_ocr(language: str):
    """
    Lazily load PaddleOCR for a given language.
    Uses runtime inspection to avoid passing unsupported args across PaddleOCR versions.
    """
    lang = language or _default_language
    with _ocr_lock:
        if lang in _ocr_models:
            return _ocr_models[lang]

        logger.info("Initializing PaddleOCR for language=%s", lang)

        # Import inside function to ensure env vars already set
        try:
            import paddle
            import paddleocr
            # Import actual class
            from paddleocr import PaddleOCR
        except Exception as exc:
            logger.exception("Failed to import paddle/paddleocr. This may indicate an incompatible wheel or missing libs.")
            raise RuntimeError(
                "Failed to import paddle/paddleocr. Ensure the installed wheel matches the platform and Python version."
            ) from exc

        # Log versions for diagnostics
        try:
            paddle_version = getattr(paddle, "__version__", "unknown")
        except Exception:
            paddle_version = "unknown"
        try:
            paddleocr_version = getattr(paddleocr, "__version__", "unknown")
        except Exception:
            paddleocr_version = "unknown"

        logger.info("Detected paddle version=%s, paddleocr version=%s", paddle_version, paddleocr_version)

        # Candidate parameters we prefer to pass if the constructor supports them
        candidate_params = {
            "lang": lang,
            "use_angle_cls": True,
            "use_gpu": False,
            # try to request single-threaded behaviour if supported
            "cpu_threads": 1,
            # try to disable mkldnn/dnnl if supported (some versions accept enable_mkldnn)
            "enable_mkldnn": False,
            # precision if supported
            "precision": "fp32",
        }

        # Inspect PaddleOCR constructor and pick only accepted args (defensive)
        try:
            sig = inspect.signature(PaddleOCR.__init__)
            ctor_params = set(sig.parameters.keys())
            # remove 'self'
            ctor_params.discard("self")
            accepted_kwargs = {k: v for k, v in candidate_params.items() if k in ctor_params}
            logger.info("PaddleOCR ctor accepted args: %s", sorted(list(accepted_kwargs.keys())))
        except Exception as exc:
            # If introspection fails, fall back to trying with minimal args
            logger.warning("Failed to inspect PaddleOCR.__init__ signature: %s. Falling back to minimal args.", exc)
            accepted_kwargs = {"lang": lang} if "lang" in candidate_params else {}

        # Finally, try to construct the PaddleOCR object with accepted args
        try:
            model = PaddleOCR(**accepted_kwargs) if accepted_kwargs else PaddleOCR()
        except Exception as exc:
            # Detect if the failure looks like a native crash / illegal instruction by message or chained C++ traces
            logger.exception("PaddleOCR initialization failed. Inspecting error...")
            # Provide a detailed runtime hint in the raised exception for the caller
            raise RuntimeError(
                "PaddleOCR initialization failed.\n"
                f"Paddle version: {paddle_version}\n"
                f"PaddleOCR version: {paddleocr_version}\n"
                "Possible causes:\n"
                "  - Incompatible paddle/paddleocr wheel for this CPU/platform (SIGILL / Illegal instruction).\n"
                "  - Native libs were preloaded with incompatible optimizations before env vars were set.\n"
                "  - Model-specific IR fusion pass triggered an illegal instruction for this build.\n\n"
                "Suggested actions:\n"
                "  1) Ensure env vars in server.py are at file top before any heavy imports (they are).\n"
                "  2) Check host CPU flags: run `lscpu | grep -i avx` on the host.\n"
                "  3) If the host lacks AVX, use a no-AVX paddle wheel or build from source.\n"
                "  4) If problem persists only for particular language/model, try a different language model (e.g., 'en') to isolate model-specific fuse issues.\n"
                "  5) As a fallback, try an earlier paddle/paddleocr combination known to be stable on this host (e.g., paddlepaddle 2.5.x + paddleocr 2.x).\n\n"
                f"Original error: {exc}"
            ) from exc

        _ocr_models[lang] = model
        logger.info("PaddleOCR initialized successfully for language=%s", lang)
        return model


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
        try:
            # PaddleOCR.ocr API changed in 3.x: the instance may expose ocr(image) or a pipeline API.
            # We attempt to call the high-level .ocr(...) and if it fails, raise a helpful error.
            results = ocr.ocr(image)
        except TypeError:
            # Older/newer APIs accept multiple params; try a more explicit call
            try:
                results = ocr.ocr(image, det=True, rec=True)
            except Exception as exc:
                logger.exception("OCR runtime error on image page with fallback signature.")
                raise HTTPException(status_code=500, detail=f"OCR runtime error: {exc}") from exc
        except Exception as exc:
            logger.exception("OCR runtime error on image page.")
            raise HTTPException(status_code=500, detail=f"OCR runtime error: {exc}") from exc

        # Normalize results from various PaddleOCR versions:
        # Some versions return list of pages -> list of items; others return list directly.
        # We'll flatten conservatively.
        flattened = []
        if not results:
            flattened = []
        else:
            # If results is a list of lists of items (page -> items), flatten two levels
            if isinstance(results, list) and results and isinstance(results[0], list) and results[0] and isinstance(results[0][0], list):
                for page in results:
                    for item in page:
                        flattened.append(item)
            else:
                flattened = results

        for item in flattened:
            # item commonly: [box, (text, confidence)] or [box, text] depending on version
            try:
                if isinstance(item, (list, tuple)) and len(item) >= 2:
                    # If second element is tuple/list with text at index 0
                    second = item[1]
                    if isinstance(second, (list, tuple)) and len(second) >= 1:
                        text = second[0]
                    else:
                        text = second
                    if text:
                        lines.append(text)
            except Exception:
                # ignore malformed items
                continue

    if not lines:
        raise HTTPException(status_code=422, detail="OCR did not return any text")

    return lines


@app.get("/healthz")
def healthz() -> dict:
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