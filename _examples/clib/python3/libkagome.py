from __future__ import annotations

import ctypes
import platform
from pathlib import Path
from typing import List


# ---------------------------------------------------------------------------
# Platform helpers
# ---------------------------------------------------------------------------


def _shared_library_name() -> str:
    match platform.system():
        case "Windows":
            return "libkagome.dll"
        case "Darwin":
            return "libkagome.dylib"
        case "Linux":
            return "libkagome.so"
        case _:
            raise RuntimeError("Unsupported OS")


def _shared_library_path() -> Path:
    base = Path(__file__).resolve().parent.parent
    path = (
        base / "bin" / _shared_library_name()
    )  # like filepath.Join(base, "bin", libname)
    if not path.exists():
        raise FileNotFoundError(f"Shared library not found: {path}")
    return path


# ---------------------------------------------------------------------------
# ctypes struct definitions (must match C layout exactly)
# ---------------------------------------------------------------------------


class _Token(ctypes.Structure):
    _fields_ = [
        ("surface", ctypes.c_char_p),
        ("pos1", ctypes.c_char_p),
        ("pos2", ctypes.c_char_p),
        ("pos3", ctypes.c_char_p),
        ("pos4", ctypes.c_char_p),
        ("base_form", ctypes.c_char_p),
        ("conj_type", ctypes.c_char_p),
        ("conj_form", ctypes.c_char_p),
        ("reading", ctypes.c_char_p),
        ("pronunciation", ctypes.c_char_p),
        ("start", ctypes.c_int),
        ("end", ctypes.c_int),
    ]


class _TokenArray(ctypes.Structure):
    _fields_ = [
        ("tokens", ctypes.POINTER(_Token)),
        ("length", ctypes.c_int),
    ]


# ---------------------------------------------------------------------------
# Public Python-side Token
# ---------------------------------------------------------------------------


class Token:
    """Pure Python representation of a Kagome token."""

    def __init__(
        self,
        surface: str,
        pos: list[str],
        base_form: str,
        conj_type: str,
        conj_form: str,
        reading: str,
        pronunciation: str,
        start: int,
        end: int,
    ) -> None:
        self.surface = surface
        self.pos = pos
        self.base_form = base_form
        self.conj_type = conj_type
        self.conj_form = conj_form
        self.reading = reading
        self.pronunciation = pronunciation
        self.start = start
        self.end = end

    def __str__(self) -> str:
        return (
            f"surface={self.surface}, "
            f"pos={self.pos}, "
            f"base_form={self.base_form}, "
            f"conj_type={self.conj_type}, "
            f"conj_form={self.conj_form}, "
            f"reading={self.reading}, "
            f"pronunciation={self.pronunciation}, "
            f"start={self.start}, end={self.end}"
        )


# ---------------------------------------------------------------------------
# Kagome wrapper
# ---------------------------------------------------------------------------


class Kagome:
    """Python-friendly wrapper around Kagome C ABI."""

    def __init__(self) -> None:
        lib_path = _shared_library_path()
        self._lib = ctypes.CDLL(str(lib_path))

        self._lib.KagomeInit.argtypes = []
        self._lib.KagomeInit.restype = ctypes.c_void_p

        self._lib.KagomeTokenizeStruct.argtypes = [
            ctypes.c_void_p,
            ctypes.c_char_p,
        ]
        self._lib.KagomeTokenizeStruct.restype = ctypes.POINTER(_TokenArray)

        self._lib.KagomeFreeTokenArray.argtypes = [ctypes.POINTER(_TokenArray)]
        self._lib.KagomeFreeTokenArray.restype = None

        self._handle = self._lib.KagomeInit()
        if not self._handle:
            raise RuntimeError("Failed to initialize Kagome")

    # ------------------------------------------------------------------

    def tokenize(self, text: str) -> List[Token]:
        """Equivalent of kagome.Tokenize."""
        arr_p = self._lib.KagomeTokenizeStruct(self._handle, text.encode("utf-8"))
        if not arr_p:
            raise RuntimeError("Tokenization failed")

        tokens: list[Token] = []
        try:
            arr = arr_p.contents
            for i in range(arr.length):
                t = arr.tokens[i]
                tokens.append(
                    Token(
                        surface=t.surface.decode("utf-8"),
                        pos=[
                            t.pos1.decode("utf-8"),
                            t.pos2.decode("utf-8"),
                            t.pos3.decode("utf-8"),
                            t.pos4.decode("utf-8"),
                        ],
                        base_form=t.base_form.decode("utf-8"),
                        conj_type=t.conj_type.decode("utf-8"),
                        conj_form=t.conj_form.decode("utf-8"),
                        reading=t.reading.decode("utf-8"),
                        pronunciation=t.pronunciation.decode("utf-8"),
                        start=t.start,
                        end=t.end,
                    )
                )
        finally:
            self._lib.KagomeFreeTokenArray(arr_p)

        return tokens

    # ------------------------------------------------------------------

    def wakati(self, text: str) -> List[str]:
        """Equivalent of kagome.Wakati."""
        return [t.surface for t in self.tokenize(text)]
