
<?php

// Detect shared library name by platform
$libnames = [
	'libkagome.so',
	'libkagome.dylib',
	'libkagome.dll',
];
$libpath = null;
foreach ($libnames as $name) {
	$try = __DIR__ . '/../bin/' . $name;
	if (file_exists($try)) {
		$libpath = $try;
		break;
	}
}
if (!$libpath) {
	fwrite(STDERR, "libkagome shared library not found\n");
	exit(1);
}

// Load FFI
$ffi = FFI::cdef('
typedef struct {
	char* surface;
	char* pos;
	int start;
	int end;
} Token;
typedef struct {
	Token* tokens;
	int length;
} TokenArray;
typedef unsigned long uintptr_t;
TokenArray* KagomeTokenizeStruct(uintptr_t handle, char* input);
void KagomeFreeTokenArray(TokenArray* arr);
uintptr_t KagomeInit(char* dictPath);
', $libpath);

// Prepare input
$text = "すもももももももものうち";

// Go: uintptr_t KagomeInit(char* dictPath);
$handle = $ffi->KagomeInit($ffi->new('char[1]', false));
if ($handle == 0) {
	fwrite(STDERR, "Failed to initialize Kagome tokenizer\n");
	exit(1);
}

// Go: TokenArray* KagomeTokenizeStruct(uintptr_t handle, char* input);
$cstr = $ffi->new('char[' . (strlen($text) + 1) . ']', false);
FFI::memcpy($cstr, $text, strlen($text));
$arr_p = $ffi->KagomeTokenizeStruct($handle, $cstr);
if ($arr_p == null) {
	fwrite(STDERR, "tokenize failed\n");
	exit(1);
}
$arr = $arr_p[0];

$expect = [
	"surface=すもも, pos=名詞,一般,*,*, start=0, end=3",
	"surface=も, pos=助詞,係助詞,*,*, start=3, end=4",
	"surface=もも, pos=名詞,一般,*,*, start=4, end=6",
	"surface=も, pos=助詞,係助詞,*,*, start=6, end=7",
	"surface=もも, pos=名詞,一般,*,*, start=7, end=9",
	"surface=の, pos=助詞,連体化,*,*, start=9, end=10",
	"surface=うち, pos=名詞,非自立,副詞可能,*, start=10, end=12",
];
$actual = [];

if ($arr->tokens != null && $arr->length > 0) {
	for ($i = 0; $i < $arr->length; $i++) {
		$token = $arr->tokens[$i];
		$line = sprintf(
			"surface=%s, pos=%s, start=%d, end=%d",
			FFI::string($token->surface),
			FFI::string($token->pos),
			$token->start,
			$token->end
		);
		echo $line . "\n";
		$actual[] = $line;
	}
	$ffi->KagomeFreeTokenArray($arr_p);
}

if ($actual === $expect) {
	echo "PASS\n";
	exit(0);
} else {
	echo "FAIL\nexpect:\n";
	foreach ($expect as $line) echo $line . "\n";
	echo "actual:\n";
	foreach ($actual as $line) echo $line . "\n";
	exit(1);
}
