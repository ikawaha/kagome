from __future__ import annotations

import ctypes
import platform
from pathlib import Path
from typing import List


"""
Python wrapper for Kagome tokenizer.

Loads the Kagome shared library and provides Python-friendly APIs:

- Kagome.tokenize(text) -> list[Token]  # Full morphological analysis
- Kagome.wakati(text) -> list[str]      # Surface forms only

The library calls kagome_* C wrapper functions, which provide a stable
interface to the Go implementation.
"""


# ---------------------------------------------------------------------------
# Platform helpers
# ---------------------------------------------------------------------------


def _shared_library_name() -> str:
    """
    Get library filename for current platform.

    Returns:
        "libkagome.dll" (Windows), "libkagome.dylib" (macOS), or "libkagome.so" (Linux)
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
    Find the shared library.

    Looks for <project_root>/bin/libkagome.(so|dll|dylib)

    Returns:
        Absolute path to the library

    Raises:
        FileNotFoundError: If library not found
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
    """Internal: C Token struct for ctypes (must match C layout exactly)."""

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
    """Internal: C TokenArray struct for ctypes."""

    _fields_ = [
        ("tokens", ctypes.POINTER(_Token)),
        ("length", ctypes.c_int),
    ]


# ---------------------------------------------------------------------------
# Public Python-side Token
# ---------------------------------------------------------------------------


class Token:
    """
    One morphological token from Kagome.

    Attributes:
        surface: Surface form (表層形)
        pos: Part-of-speech [major, middle, small, detail] (品詞階層)
        base_form: Base form (原形)
        conj_type: Conjugation type (活用型)
        conj_form: Conjugation form (活用形)
        reading: Katakana reading (読み)
        pronunciation: Pronunciation in katakana (発音)
        start: Start position in input (開始位置)
        end: End position in input (終了位置)
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
        """String representation for display and testing."""
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
    Python wrapper for Kagome tokenizer.

    Loads the shared library, calls kagome_* wrapper functions,
    and converts results to Python objects.

    Create once and reuse for multiple tokenizations.
    """

    def __init__(self) -> None:
        """
        Load library and create tokenizer.

        Raises:
            RuntimeError: If initialization fails
        """
        lib_path = _shared_library_path()
        self._lib = ctypes.CDLL(str(lib_path))

        # Define C function signatures
        self._lib.kagome_init.argtypes = []
        self._lib.kagome_init.restype = ctypes.c_void_p

        self._lib.kagome_tokenize.argtypes = [
            ctypes.c_void_p,
            ctypes.c_char_p,
        ]
        self._lib.kagome_tokenize.restype = ctypes.POINTER(_TokenArray)

        self._lib.kagome_free_token_array.argtypes = [ctypes.POINTER(_TokenArray)]
        self._lib.kagome_free_token_array.restype = None

        self._lib.kagome_destroy.argtypes = [ctypes.c_void_p]
        self._lib.kagome_destroy.restype = None

        self._handle = self._lib.kagome_init()
        if not self._handle:
            raise RuntimeError("Failed to initialize Kagome")

    def __del__(self) -> None:
        """Clean up tokenizer when object is deleted."""
        if hasattr(self, "_handle") and self._handle:
            self._lib.kagome_destroy(self._handle)
            self._handle = None

    # ------------------------------------------------------------------

    def tokenize(self, text: str) -> List[Token]:
        """
        Tokenize Japanese text.

        Args:
            text: Japanese text to tokenize

        Returns:
            List of Token objects with full morphological info

        Raises:
            RuntimeError: If tokenization fails
        """
        arr_p = self._lib.kagome_tokenize(self._handle, text.encode("utf-8"))
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
            self._lib.kagome_free_token_array(arr_p)

        return tokens

    # ------------------------------------------------------------------

    def wakati(self, text: str) -> List[str]:
        """
        Tokenize and return surface forms only (wakati-gaki 分かち書き).

        Args:
            text: Japanese text to tokenize

        Returns:
            List of surface forms (without morphological info)
        """
        return [t.surface for t in self.tokenize(text)]
