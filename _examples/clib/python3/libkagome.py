from __future__ import annotations

import ctypes
import platform
from pathlib import Path
from typing import List


"""
libkagome.py

Python wrapper for Kagome tokenizer via C ABI.

This module loads a shared library built from the Go implementation of Kagome
and exposes Python-friendly APIs equivalent to:

- kagome.Tokenize(input string) []Token
- kagome.Wakati(input string) []string
"""


# ---------------------------------------------------------------------------
# Platform helpers
# ---------------------------------------------------------------------------


def _shared_library_name() -> str:
    """
    Return platform-specific shared library filename.

    Returns:
        str:
            - Windows: libkagome.dll
            - macOS:   libkagome.dylib
            - Linux:   libkagome.so
    """
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
    """
    Resolve absolute path to the Kagome shared library.

    The library is expected at:

        <project_root>/bin/libkagome.(so|dll|dylib)

    Note:
        Path joining via `/` is OS-independent in pathlib,
        similar to Go's filepath.Join.

    Returns:
        Path: Absolute path to the shared library.

    Raises:
        FileNotFoundError: If the shared library does not exist.
    """
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
    """
    ctypes representation of Kagome Token struct (C ABI).

    This structure must match the C struct layout exactly.
    """

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
    """
    ctypes representation of array returned from KagomeTokenizeStruct.
    """

    _fields_ = [
        ("tokens", ctypes.POINTER(_Token)),
        ("length", ctypes.c_int),
    ]


# ---------------------------------------------------------------------------
# Public Python-side Token
# ---------------------------------------------------------------------------


class Token:
    """
    Python representation of a Kagome token.

    Each field corresponds to Kagome's dictionary features.

    Attributes:
        surface (str):
            Surface form (表層形)

        pos (list[str]):
            Part-of-speech hierarchy (品詞階層)
            - pos[0]: Major category (品詞大分類)
            - pos[1]: Subcategory 1 (品詞中分類)
            - pos[2]: Subcategory 2 (品詞小分類)
            - pos[3]: Subcategory 3 / detail (品詞再分類)

        base_form (str):
            Base / dictionary form (原型・基本形)

        conj_type (str):
            Conjugation type (活用型)
            e.g. 五段・カ行促音便

        conj_form (str):
            Conjugation form (活用形)
            e.g. 連用タ接続

        reading (str):
            Reading in katakana (読み)
            e.g. 公園 -> コウエン

        pronunciation (str):
            Pronunciation (発音)
            e.g. 公園 -> コーエン

        start (int):
            Start position in input string (開始位置)

        end (int):
            End position in input string (終了位置)
    """

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
        """Human-readable representation used in examples and tests."""
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
    """
    Python-friendly wrapper around Kagome C ABI.

    This class manages:
    - Loading the shared library
    - Initializing the Kagome tokenizer
    - Converting C structs into Python objects

    Intended to be long-lived and reused.
    """

    def __init__(self) -> None:
        """
        Load Kagome shared library and initialize tokenizer.

        Raises:
            RuntimeError: If initialization fails.
        """
        lib_path = _shared_library_path()
        self._lib = ctypes.CDLL(str(lib_path))

        # C function signatures
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
        """
        Tokenize input text (equivalent of kagome.Tokenize).

        Args:
            text (str): Input text to tokenize.

        Returns:
            list[Token]: List of Token objects with full morphological info.

        Raises:
            RuntimeError: If tokenization fails.
        """
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
            # Always free C-allocated memory
            self._lib.KagomeFreeTokenArray(arr_p)

        return tokens

    # ------------------------------------------------------------------

    def wakati(self, text: str) -> List[str]:
        """
        Perform wakati-gaki (分かち書き).

        Equivalent of kagome.Wakati.

        Args:
            text (str): Input text.

        Returns:
            list[str]: List of surface forms only.
        """
        return [t.surface for t in self.tokenize(text)]
