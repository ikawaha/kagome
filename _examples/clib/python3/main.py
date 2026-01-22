from __future__ import annotations

import ctypes
import platform
import sys
from pathlib import Path


# ---------------------------------------------------------------------------
# Platform helpers
# ---------------------------------------------------------------------------


def shared_library_name() -> str:
    """Return platform-specific shared library name."""
    match platform.system():
        case "Windows":
            return "libkagome.dll"
        case "Darwin":
            return "libkagome.dylib"
        case "Linux":
            return "libkagome.so"
        case _:
            raise RuntimeError("Unsupported OS")


def shared_library_path() -> Path:
    """Resolve shared library path under ../bin."""
    path = Path(__file__).resolve().parent.parent / "bin" / shared_library_name()
    if not path.exists():
        raise FileNotFoundError(f"Shared library not found: {path}")
    return path


# ---------------------------------------------------------------------------
# ctypes struct definitions (must match C layout exactly)
# ---------------------------------------------------------------------------


class Token(ctypes.Structure):
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


class TokenArray(ctypes.Structure):
    _fields_ = [
        ("tokens", ctypes.POINTER(Token)),
        ("length", ctypes.c_int),
    ]


# ---------------------------------------------------------------------------
# Kagome FFI loader
# ---------------------------------------------------------------------------


class KagomeFFI:
    """Thin Python wrapper around Kagome C ABI."""

    def __init__(self, lib_path: Path) -> None:
        print(f"Loading library: {lib_path}")
        self._lib = ctypes.CDLL(str(lib_path))

        # function signatures
        self._lib.KagomeInit.argtypes = []
        self._lib.KagomeInit.restype = ctypes.c_void_p

        self._lib.KagomeTokenizeStruct.argtypes = [
            ctypes.c_void_p,
            ctypes.c_char_p,
        ]
        self._lib.KagomeTokenizeStruct.restype = ctypes.POINTER(TokenArray)

        self._lib.KagomeFreeTokenArray.argtypes = [ctypes.POINTER(TokenArray)]
        self._lib.KagomeFreeTokenArray.restype = None

    def init(self) -> ctypes.c_void_p:
        handle = self._lib.KagomeInit()
        if not handle:
            raise RuntimeError("Failed to initialize Kagome tokenizer")
        return handle

    def tokenize(self, handle: ctypes.c_void_p, text: bytes) -> list[str]:
        arr_p = self._lib.KagomeTokenizeStruct(handle, text)
        if not arr_p:
            raise RuntimeError("Tokenization failed")

        results: list[str] = []
        try:
            arr = arr_p.contents
            for i in range(arr.length):
                tok = arr.tokens[i]
                results.append(format_token(tok))
        finally:
            self._lib.KagomeFreeTokenArray(arr_p)

        return results


# ---------------------------------------------------------------------------
# Formatting helpers
# ---------------------------------------------------------------------------


def decode(ptr: ctypes.c_char_p) -> str:
    return ptr.decode("utf-8")


def format_token(token: Token) -> str:
    pos = [
        decode(token.pos1),
        decode(token.pos2),
        decode(token.pos3),
        decode(token.pos4),
    ]
    return (
        f"surface={decode(token.surface)}, "
        f"pos={pos}, "
        f"base_form={decode(token.base_form)}, "
        f"conj_type={decode(token.conj_type)}, "
        f"conj_form={decode(token.conj_form)}, "
        f"reading={decode(token.reading)}, "
        f"pronunciation={decode(token.pronunciation)}, "
        f"start={token.start}, end={token.end}"
    )


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------


def main() -> int:
    ffi = KagomeFFI(shared_library_path())
    handle = ffi.init()

    text = "すもももももももものうち".encode("utf-8")

    expect = [
        "surface=すもも, pos=['名詞', '一般', '*', '*'], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3",
        "surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4",
        "surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6",
        "surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7",
        "surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9",
        "surface=の, pos=['助詞', '連体化', '*', '*'], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10",
        "surface=うち, pos=['名詞', '非自立', '副詞可能', '*'], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12",
    ]

    actual = ffi.tokenize(handle, text)

    for line in actual:
        print(line)

    if actual == expect:
        print("PASS")
        return 0

    print("FAIL\nexpect:")
    print("\n".join(expect))
    print("actual:")
    print("\n".join(actual))
    return 1


if __name__ == "__main__":
    sys.exit(main())
