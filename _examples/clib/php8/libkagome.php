<?php
declare(strict_types=1);

/**
 * Kagome PHP FFI wrapper.
 *
 * This class loads the Kagome C ABI shared library and provides
 * PHP-friendly methods such as tokenize() and wakati().
 */
final class Kagome
{
    private FFI $ffi;
    private mixed $handle;

    public function __construct()
    {
        $libPath = self::resolveLibraryPath();
        $this->ffi = FFI::cdef(self::cDefinitions(), $libPath);

        $this->handle = $this->ffi->KagomeInit();
        if ($this->handle === null) {
            throw new RuntimeException('Failed to initialize Kagome tokenizer');
        }
    }

    public function __destruct()
    {
        if (isset($this->handle)) {
            $this->ffi->KagomeDestroy($this->handle);
        }
    }

    /**
     * Tokenize input text.
     *
     * Equivalent to kagome.Tokenize in Go.
     *
     * @param string $text UTF-8 encoded Japanese text
     * @return Token[] List of tokens
     */
    public function tokenize(string $text): array
    {
        $cstr = $this->ffi->new('char[' . (strlen($text) + 1) . ']', false);
        FFI::memcpy($cstr, $text, strlen($text));

        $arrPtr = $this->ffi->KagomeTokenizeStruct($this->handle, $cstr);
        if ($arrPtr === null) {
            throw new RuntimeException('Tokenization failed');
        }

        $tokens = [];

        try {
            $arr = $arrPtr[0];
            for ($i = 0; $i < $arr->length; $i++) {
                $tokens[] = Token::fromC($arr->tokens[$i]);
            }
        } finally {
            $this->ffi->KagomeFreeTokenArray($arrPtr);
        }

        return $tokens;
    }

    /**
     * Wakati (surface-only tokenization).
     *
     * Equivalent to kagome.Wakati in Go.
     *
     * @param string $text
     * @return string[] List of token surfaces
     */
    public function wakati(string $text): array
    {
        return array_map(
            static fn (Token $t) => $t->surface,
            $this->tokenize($text)
        );
    }

    // ---------------------------------------------------------------------

    private static function resolveLibraryPath(): string
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
     * Raw C ABI definitions for FFI.
     *
     * IMPORTANT:
     * - Do not add comments here.
     * - Keep this in sync with the C header.
     */
    private static function cDefinitions(): string
    {
        return <<<CDEF
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
    }
}

/**
 * Immutable PHP representation of a Kagome token.
 *
 * Each field corresponds to Kagome dictionary features.
 */
final class Token
{
    /** Surface form (表層形) */
    public string $surface;

    /**
     * Part-of-speech hierarchy (品詞階層)
     *
     * [0] Major class (大分類)
     * [1] Middle class (中分類)
     * [2] Small class (小分類)
     * [3] Fine class / subcategory (再分類)
     *
     * @var string[]
     */
    public array $pos;

    /** Base form / dictionary form (原形・基本形) */
    public string $base_form;

    /** Conjugation type (活用型, e.g. 五段・カ行促音便) */
    public string $conj_type;

    /** Conjugation form (活用形, e.g. 連用タ接続) */
    public string $conj_form;

    /** Reading in katakana (読み, e.g. 公園 → コウエン) */
    public string $reading;

    /** Pronunciation (発音, e.g. 公園 → コーエン) */
    public string $pronunciation;

    /** Start byte position in input text (開始位置) */
    public int $start;

    /** End byte position in input text (終了位置) */
    public int $end;

    private function __construct() {}

    /**
     * Create a Token from a C Token struct.
     *
     * @internal
     */
    public static function fromC(object $t): self
    {
        $self = new self();
        $self->surface = FFI::string($t->surface);
        $self->pos = [
            FFI::string($t->pos1),
            FFI::string($t->pos2),
            FFI::string($t->pos3),
            FFI::string($t->pos4),
        ];
        $self->base_form = FFI::string($t->base_form);
        $self->conj_type = FFI::string($t->conj_type);
        $self->conj_form = FFI::string($t->conj_form);
        $self->reading = FFI::string($t->reading);
        $self->pronunciation = FFI::string($t->pronunciation);
        $self->start = $t->start;
        $self->end = $t->end;

        return $self;
    }

    public function __toString(): string
    {
        return sprintf(
            'surface=%s, pos=[%s], base_form=%s, conj_type=%s, conj_form=%s, reading=%s, pronunciation=%s, start=%d, end=%d',
            $this->surface,
            implode(', ', $this->pos),
            $this->base_form,
            $this->conj_type,
            $this->conj_form,
            $this->reading,
            $this->pronunciation,
            $this->start,
            $this->end
        );
    }
}
