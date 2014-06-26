package trie

type Trie interface {
	Search(string) (id int, ok bool)
	PrefixSearch(string) (keyword string, id int, ok bool)
	CommonPrefixSearch(string) (keywords []string, ids []int)
}
