import ctypes
import platform
import os
import sys

# Select library name by OS
system = platform.system()
if system == "Windows":
    libname = "libkagome.dll"
elif system == "Darwin":
    libname = "libkagome.dylib"
else:
    libname = "libkagome.so"

# ./.. /bin/libkagome(_linux.so|_win.dll|_mac.dylib)
libpath = os.path.abspath(
    os.path.join(os.path.dirname(__file__), os.path.pardir, "bin", libname)
)
print(f"Loading library: {libpath}")

# lib = ctypes.cdll.LoadLibrary(libpath)
lib = ctypes.CDLL(os.path.abspath(libpath), mode=ctypes.RTLD_GLOBAL)
lib.KagomeInit.argtypes = []
lib.KagomeInit.restype = ctypes.c_void_p


# Define Token and TokenArray struct for ctypes
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


# Set argument and return types for Go shared library functions (pointer version)
lib.KagomeTokenizeStruct.argtypes = [ctypes.c_size_t, ctypes.c_char_p]
lib.KagomeTokenizeStruct.restype = ctypes.POINTER(TokenArray)
lib.KagomeFreeTokenArray.argtypes = [ctypes.POINTER(TokenArray)]
lib.KagomeFreeTokenArray.restype = None

# Initialize tokenizer and get handle
handle = lib.KagomeInit()
if handle == 0:
    raise RuntimeError("Failed to initialize Kagome tokenizer")


# Tokenize text using struct array
text = "すもももももももものうち".encode("utf-8")


# --- Testable output ---
expect = [
    "surface=すもも, pos=['名詞', '一般', '*', '*'], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3",
    "surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4",
    "surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6",
    "surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7",
    "surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9",
    "surface=の, pos=['助詞', '連体化', '*', '*'], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10",
    "surface=うち, pos=['名詞', '非自立', '副詞可能', '*'], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12",
]

actual = []

# Call KagomeTokenizeStruct and handle pointer result
arr_p = lib.KagomeTokenizeStruct(handle, text)
if not arr_p:
    raise RuntimeError("tokenize failed")

arr = arr_p.contents
if arr.tokens and arr.length > 0:
    for i in range(arr.length):
        token = arr.tokens[i]
        surface = token.surface.decode("utf-8")
        pos_arr = [
            token.pos1.decode("utf-8"),
            token.pos2.decode("utf-8"),
            token.pos3.decode("utf-8"),
            token.pos4.decode("utf-8"),
        ]
        line = f"surface={surface}, pos={pos_arr}, base_form={token.base_form.decode('utf-8')}, conj_type={token.conj_type.decode('utf-8')}, conj_form={token.conj_form.decode('utf-8')}, reading={token.reading.decode('utf-8')}, pronunciation={token.pronunciation.decode('utf-8')}, start={token.start}, end={token.end}"
        print(line)
        actual.append(line)
    lib.KagomeFreeTokenArray(arr_p)

if actual == expect:
    print("PASS")
    sys.exit(0)
else:
    print("FAIL")
    print("expect:")
    for line in expect:
        print(line)
    print("actual:")
    for line in actual:
        print(line)
    sys.exit(1)
