<?php
declare(strict_types=1);

/**
 * Resolve shared library path by OS.
 */
function resolveLibraryPath(): string
{
    $libName = match (PHP_OS_FAMILY) {
        'Windows' => 'libkagome.dll',
        'Darwin'  => 'libkagome.dylib',
        'Linux'   => 'libkagome.so',
        default   => throw new RuntimeException('Unsupported OS family'),
    };

    $path = realpath(__DIR__ . '/../bin/' . $libName);
    if ($path === false || !is_file($path)) {
        throw new RuntimeException("Shared library not found: {$libName}");
    }

    return $path;
}

/**
 * C definitions for Kagome C ABI.
 */
const KAGOME_CDEF = <<<CDEF
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
void KagomeDestroy(void* handle);
CDEF;

try {
    $libPath = resolveLibraryPath();
    echo "Loading library: {$libPath}" . PHP_EOL;

    $ffi = FFI::cdef(KAGOME_CDEF, $libPath);

    // Initialize tokenizer
    $handle = $ffi->KagomeInit();
    if ($handle === null) {
        throw new RuntimeException('Failed to initialize Kagome tokenizer');
    }

    $text = 'すもももももももものうち';

    // Prepare C string (null-terminated)
    $cstr = $ffi->new('char[' . (strlen($text) + 1) . ']', false);
    FFI::memcpy($cstr, $text, strlen($text));

    $arrPtr = $ffi->KagomeTokenizeStruct($handle, $cstr);
    if ($arrPtr === null) {
        throw new RuntimeException('Tokenization failed');
    }

    $expected = [
        "surface=すもも, pos=[名詞, 一般, *, *], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3",
        "surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4",
        "surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6",
        "surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7",
        "surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9",
        "surface=の, pos=[助詞, 連体化, *, *], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10",
        "surface=うち, pos=[名詞, 非自立, 副詞可能, *], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12",
    ];

    $actual = [];

    try {
        $arr = $arrPtr[0];
        for ($i = 0; $i < $arr->length; $i++) {
            $t = $arr->tokens[$i];

            $pos = implode(', ', [
                FFI::string($t->pos1),
                FFI::string($t->pos2),
                FFI::string($t->pos3),
                FFI::string($t->pos4),
            ]);

            $line = sprintf(
                'surface=%s, pos=[%s], base_form=%s, conj_type=%s, conj_form=%s, reading=%s, pronunciation=%s, start=%d, end=%d',
                FFI::string($t->surface),
                $pos,
                FFI::string($t->base_form),
                FFI::string($t->conj_type),
                FFI::string($t->conj_form),
                FFI::string($t->reading),
                FFI::string($t->pronunciation),
                $t->start,
                $t->end
            );

            echo $line . PHP_EOL;
            $actual[] = $line;
        }
    } finally {
        // Always free memory allocated by Go
        $ffi->KagomeFreeTokenArray($arrPtr);
    }

	// Success
    if ($actual === $expected) {
        echo "PASS" . PHP_EOL;
        $ffi->KagomeDestroy($handle);
        exit(0);
    }

	// Failure
    echo "FAIL" . PHP_EOL . "expect:" . PHP_EOL;
    foreach ($expected as $line) {
        echo $line . PHP_EOL;
    }

    echo "actual:" . PHP_EOL;
    foreach ($actual as $line) {
        echo $line . PHP_EOL;
    }

    $ffi->KagomeDestroy($handle);
    exit(1);

} catch (Throwable $e) {
    fwrite(STDERR, $e->getMessage() . PHP_EOL);
    exit(1);
}
