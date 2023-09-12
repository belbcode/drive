package filesystem

import "strings"

func (t *Tree) Search(searchQuery string) []Node {
	var results []Node
	contains := func(n Node) {
		if strings.Contains(n.Info.Name(), searchQuery) {
			results = append(results, n)
		}

	}
	t.Traverse(contains)
	return results
}
