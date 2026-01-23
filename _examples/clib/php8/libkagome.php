<?php
declare(strict_types=1);

/**
 * PHP wrapper for Kagome tokenizer.
 *
 * Loads the shared library and provides PHP-friendly methods:
 * - tokenize(): Full morphological analysis
 * - wakati(): Surface forms only
 */
final class Kagome
{
    private FFI $ffi;
    private mixed $handle;

    public function __construct()
    {
        $libPath = self::resolveLibraryPath();
        $this->ffi = FFI::cdef(self::cDefinitions(), $libPath);

        $this->handle = $this->ffi->kagome_init();
        if ($this->handle === null) {
            throw new RuntimeException('Failed to initialize Kagome tokenizer');
        }
    }

    public function __destruct()
    {
        if (isset($this->handle)) {
            $this->ffi->kagome_destroy($this->handle);
        }
    }

    /**
     * Tokenize Japanese text.
     *
     * @param string $text Japanese text (UTF-8)
     * @return Token[] Tokens with full morphological info
     */
    public function tokenize(string $text): array
    {
        $cstr = $this->ffi->new('char[' . (strlen($text) + 1) . ']', false);
        FFI::memcpy($cstr, $text, strlen($text));

        $arrPtr = $this->ffi->kagome_tokenize($this->handle, $cstr);
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
            $this->ffi->kagome_free_token_array($arrPtr);
        }

        return $tokens;
    }

    /**
     * Tokenize and return surface forms only (wakati-gaki).
     *
     * @param string $text Japanese text (UTF-8)
     * @return string[] Surface forms only
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
     * C function declarations for FFI.
     *
     * Must match kagome_wrapper.h exactly.
     * Do not add comments inside the CDEF string.
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

void* kagome_init(void);
void kagome_destroy(void* handle);
TokenArray* kagome_tokenize(void* handle, const char* input);
void kagome_free_token_array(TokenArray* arr);
CDEF;
    }
}

/**
 * One morphological token from Kagome.
 */
final class Token
{
    /** Surface form (表層形) */
    public string $surface;

    /** Part-of-speech [major, middle, small, detail] (品詞階層) */
    public array $pos;

    /** Base form (原形) */
    public string $base_form;

    /** Conjugation type (活用型) */
    public string $conj_type;

    /** Conjugation form (活用形) */
    public string $conj_form;

    /** Katakana reading (読み) */
    public string $reading;

    /** Pronunciation (発音) */
    public string $pronunciation;

    /** Start position in input */
    public int $start;

    /** End position in input */
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
