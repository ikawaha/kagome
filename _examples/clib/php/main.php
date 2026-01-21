<?php

// Detect shared library name by platform
$libnames = [
	'libkagome.so',
	'libkagome.dylib',
	'libkagome.dll',
];
$libpath = null;
foreach ($libnames as $name) {
    $try = realpath(__DIR__ . '/../bin/' . $name);
    if ($try !== false && file_exists($try)) {
        $libpath = $try;
        break;
    }
}

if (!$libpath) {
	fwrite(STDERR, "libkagome shared library not found" . PHP_EOL);
	exit(1);
}
echo "Loading library: {$libpath}" . PHP_EOL;

// Load FFI
 $ffi = FFI::cdef('
typedef struct {
	char* surface;
	char* pos1;
	char* pos2;
	char* pos3;
	char* pos4;
	char* base_form;
	char* conj_type;
	char* conj_form;
	char* reading;
	char* pronunciation;
	int start;
	int end;
} Token;
typedef struct {
	Token* tokens;
	int length;
} TokenArray;
TokenArray* KagomeTokenizeStruct(void* handle, char* input);
void KagomeFreeTokenArray(TokenArray* arr);
void* KagomeInit(void);
', $libpath);

// Prepare input
$text = "すもももももももものうち";

// Go: void* KagomeInit(void);
$handle = $ffi->KagomeInit();
if ($handle == null) {
	fwrite(STDERR, "Failed to initialize Kagome tokenizer" . PHP_EOL);
	exit(1);
}

// Go: TokenArray* KagomeTokenizeStruct(void* handle, char* input);
$cstr = $ffi->new('char[' . (strlen($text) + 1) . ']', false);
FFI::memcpy($cstr, $text, strlen($text));
$arr_p = $ffi->KagomeTokenizeStruct($handle, $cstr);
if ($arr_p == null) {
	fwrite(STDERR, "tokenize failed" . PHP_EOL);
	exit(1);
}
$arr = $arr_p[0];

$expect = [
	"surface=すもも, pos=[名詞, 一般, *, *], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3",
	"surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4",
	"surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6",
	"surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7",
	"surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9",
	"surface=の, pos=[助詞, 連体化, *, *], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10",
	"surface=うち, pos=[名詞, 非自立, 副詞可能, *], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12",
];
$actual = [];

if ($arr->tokens != null && $arr->length > 0) {
	try {
		for ($i = 0; $i < $arr->length; $i++) {
			$token = $arr->tokens[$i];
			$surface = FFI::string($token->surface);
			$posArr = [
				FFI::string($token->pos1), // Part-of-speech, 品詞
				FFI::string($token->pos2), // POS Subcategory1, 品詞細分類1
				FFI::string($token->pos3), // POS Subcategory2, 品詞細分類2
				FFI::string($token->pos4), // POS Subcategory3, 品詞細分類3
			];
			$line = sprintf(
				"surface=%s, pos=[%s], base_form=%s, conj_type=%s, conj_form=%s, reading=%s, pronunciation=%s, start=%d, end=%d",
				$surface,
				implode(', ', $posArr),
				FFI::string($token->base_form),
				FFI::string($token->conj_type),
				FFI::string($token->conj_form),
				FFI::string($token->reading),
				FFI::string($token->pronunciation),
				$token->start,
				$token->end
			);
			echo $line . PHP_EOL;
			$actual[] = $line;
		}
	} finally {
		$ffi->KagomeFreeTokenArray($arr_p);
		unset($arr);
		unset($arr_p);
	}
}

if ($actual === $expect) {
	echo "PASS" . PHP_EOL;
	exit(0);
} else {
	echo "FAIL" . PHP_EOL . "expect:" . PHP_EOL;
	foreach ($expect as $line) echo $line . PHP_EOL;
	echo "actual:" . PHP_EOL;
	foreach ($actual as $line) echo $line . PHP_EOL;
	exit(1);
}
