package kagome

type Trie interface {
	FindString(string) (id int, ok bool)
	CommonPrefixSearchString(string) (ids, lens []int)
}
