<?php
declare(strict_types=1);

require __DIR__ . '/libkagome.php';

$kagome = new Kagome();

$text = 'すもももももももものうち';

$expected = [
    "surface=すもも, pos=[名詞, 一般, *, *], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3",
    "surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4",
    "surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6",
    "surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7",
    "surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9",
    "surface=の, pos=[助詞, 連体化, *, *], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10",
    "surface=うち, pos=[名詞, 非自立, 副詞可能, *], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12",
];

$actual = array_map(
    static fn (Token $t) => (string)$t,
    $kagome->tokenize($text)
);

foreach ($actual as $line) {
    echo $line, PHP_EOL;
}

if ($actual === $expected) {
    echo "PASS", PHP_EOL;
    exit(0);
}

echo "FAIL", PHP_EOL;
exit(1);
