package pkg

type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
}

type Trie struct {
	Root *TrieNode
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[rune]*TrieNode),
		IsEnd:    false,
	}
}

func NewTrie() *Trie {
	return &Trie{
		Root: newTrieNode(),
	}
}

//            __________________________
//            |         root           |
//            ——————————————————————————
// 		      /           \           \
//           /             \           \
//          h               a          y
//         / \             /          /
//        e   o           r          o
//       /     \         /          /
//      l       w       e          u
//     /
//    l
//   /
//  o

func (t *Trie) Insert(word string) {
	node := t.Root
	for _, c := range word {
		if _, ok := node.Children[c]; !ok {
			node.Children[c] = newTrieNode()
		}
		node = node.Children[c]
	}
	node.IsEnd = true
	node.Children = make(map[rune]*TrieNode, 0)
}

func (t *Trie) Search(text string) bool {
	runes := []rune(text)
	for i := 0; i < len(runes); {
		node := t.Root
		for j := i; j < len(runes); j++ {
			if _, ok := node.Children[runes[j]]; ok {
				node = node.Children[runes[j]]
				if node.IsEnd {
					return true
				}
			} else {
				break
			}
		}
		i++
	}
	return false
}
