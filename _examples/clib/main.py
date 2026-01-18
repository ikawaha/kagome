
import ctypes
import platform
import os

# Select library name by OS
system = platform.system()
if system == "Windows":
	libname = "libkagome.dll"
elif system == "Darwin":
	libname = "libkagome.dylib"
else:
	libname = "libkagome.so"

libpath = os.path.join(os.path.dirname(__file__), libname)
lib = ctypes.cdll.LoadLibrary(libpath)


lib.KagomeInit.argtypes = [ctypes.c_char_p]
lib.KagomeInit.restype = ctypes.c_size_t  # uintptr_t is size_t

# Use c_void_p for Go string pointers
lib.KagomeTokenize.argtypes = [ctypes.c_size_t, ctypes.c_char_p]
lib.KagomeTokenize.restype = ctypes.c_void_p
lib.KagomeFree.argtypes = [ctypes.c_void_p]
lib.KagomeFree.restype = None

# Initialize tokenizer and get handle
handle = lib.KagomeInit(b"")
if handle == 0:
	raise RuntimeError("Failed to initialize Kagome tokenizer")

# Tokenize text using handle
# NOTE: KagomeFree must only be called with a valid pointer returned by KagomeTokenize.
# If p is None, calling KagomeFree will crash Python (pointer being freed was not allocated).
p = lib.KagomeTokenize(handle, "すもももももももものうち".encode("utf-8"))
if p:
	print(ctypes.string_at(p).decode("utf-8"))
	lib.KagomeFree(p)
