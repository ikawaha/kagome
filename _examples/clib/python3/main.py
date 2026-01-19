
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
libpath = os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir, "bin", libname))
print(f"Loading library: {libpath}")

#lib = ctypes.cdll.LoadLibrary(libpath)
lib = ctypes.CDLL(
    os.path.abspath(libpath),
    mode=ctypes.RTLD_GLOBAL
)
lib.KagomeInit.argtypes = [ctypes.c_char_p]
lib.KagomeInit.restype = ctypes.c_size_t  # uintptr_t is size_t

# Define Token and TokenArray struct for ctypes
class Token(ctypes.Structure):
	_fields_ = [
		("surface", ctypes.c_char_p),
		("pos", ctypes.c_char_p),
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
handle = lib.KagomeInit(b"")
if handle == 0:
	raise RuntimeError("Failed to initialize Kagome tokenizer")


# Tokenize text using struct array
text = "すもももももももものうち".encode("utf-8")


# --- Testable output ---
expect = [
	"surface=すもも, pos=名詞,一般,*,*, start=0, end=3",
	"surface=も, pos=助詞,係助詞,*,*, start=3, end=4",
	"surface=もも, pos=名詞,一般,*,*, start=4, end=6",
	"surface=も, pos=助詞,係助詞,*,*, start=6, end=7",
	"surface=もも, pos=名詞,一般,*,*, start=7, end=9",
	"surface=の, pos=助詞,連体化,*,*, start=9, end=10",
	"surface=うち, pos=名詞,非自立,副詞可能,*, start=10, end=12",
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
		line = f"surface={token.surface.decode('utf-8')}, pos={token.pos.decode('utf-8')}, start={token.start}, end={token.end}"
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
