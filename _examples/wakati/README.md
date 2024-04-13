# Wakati Example with Kagome

## Segmenting Japanese text into words with Kagome

In this example, we demonstrate how to segment Japanese text into words using Kagome.

- Target text data is as follows:

```text
すもももももももものうち
```

- Example output:

```shellsession
$ cd /path/to/kagome/_examples/wakati
$ go run .
----wakati---
すもも/も/もも/も/もも/の/うち
```

> __Note__ that segmentation varies depending on the dictionary used.
> In this example we use the IPA dictionary. But for searching purposes, the Uni dictionary is recommended.
>
> - [What is a Kagome dictionary?](https://github.com/ikawaha/kagome/wiki/About-the-dictionary#what-is-a-kagome-dictionary) | Wiki | kagome @ GitHub
