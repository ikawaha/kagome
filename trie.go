package kagome

// Any type implements Trie interface may be used as a dictionary.
type Trie interface {
	FindString(string) (id int, ok bool)               // search a dictionary by a keyword.
	CommonPrefixSearchString(string) (ids, lens []int) // finds keywords sharing common prefix in a dictionary.
}
