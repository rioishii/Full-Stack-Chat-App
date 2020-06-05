package indexes

import "sync"

type int64set map[int64]struct{}

func (s int64set) add(value int64) bool {
	_, ok := s[value]
	if ok {
		return false
	}
	s[value] = struct{}{}
	return true
}

func (s int64set) remove(value int64) bool {
	_, ok := s[value]
	if !ok {
		return false
	}
	delete(s, value)
	return true
}

func (s int64set) has(value int64) bool {
	_, ok := s[value]
	return ok
}

func (s int64set) all() []int64 {
	values := make([]int64, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

type trieNode struct {
	name     rune
	vals     int64set
	children map[rune]*trieNode
	parent   *trieNode
}

//Trie implements a trie data structure mapping strings to int64s
//that is safe for concurrent use.
type Trie struct {
	mx   sync.RWMutex
	Root *trieNode
	Size int
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{
		Root: &trieNode{children: make(map[rune]*trieNode)},
		Size: 0,
	}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.Size
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	t.mx.Lock()
	runes := []rune(key)
	currNode := t.Root
	for _, name := range runes {
		if currNode.children[name] == nil {
			currNode.newChild(name)
		}
		currNode = currNode.children[name]
	}
	ok := currNode.vals.add(value)
	if ok {
		t.Size++
	}
	t.mx.Unlock()
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	t.mx.RLock()
	defer t.mx.RUnlock()
	if t.Len() == 0 || len(prefix) == 0 || max == 0 {
		return nil
	}
	visited := make(map[*trieNode]bool)
	values := []int64{}
	runes := []rune(prefix)
	var dfsFind func(node *trieNode, index int)
	dfsFind = func(node *trieNode, index int) {
		if node != nil {
			if index == len(runes) || node.name == runes[index] {
				visited[node] = true
				if len(node.vals) > 0 {
					ids := node.vals.all()
					for _, id := range ids {
						if len(values) < max {
							found := searchSlice(values, id)
							if !found {
								values = append(values, id)
							}
						} else {
							break
						}
					}
				}
				if index < len(runes) {
					index++
				}
				for _, child := range node.children {
					if _, ok := visited[child]; !ok {
						dfsFind(child, index)
					}
				}
			}
		}
	}
	if n, ok := t.Root.children[runes[0]]; ok {
		dfsFind(n, 0)
	} else {
		return nil
	}
	return values
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	t.mx.Lock()
	runes := []rune(key)
	lastNode := findLastNode(t.Root, runes, value)
	var deleteNodes func(node *trieNode, index int)
	deleteNodes = func(node *trieNode, index int) {
		if node != nil && index >= 0 {
			if len(node.getChildren()) == 0 {
				parent := node.getParent()
				if parent != nil {
					parent.removeChild(runes[index])
					deleteNodes(parent, index-1)
				}
			}
		}
	}
	if lastNode != nil {
		deleteNodes(lastNode, len(runes)-1)
		t.Size--
	}
	t.mx.Unlock()
}
func (node *trieNode) newChild(name rune) {
	newNode := &trieNode{
		name:     name,
		vals:     make(int64set),
		children: make(map[rune]*trieNode),
		parent:   node,
	}
	node.children[name] = newNode
}
func (node *trieNode) removeChild(name rune) {
	delete(node.children, name)
}
func findLastNode(node *trieNode, runes []rune, value int64) *trieNode {
	if node == nil {
		return nil
	}
	if len(runes) == 0 {
		ok := node.vals.remove(value)
		if ok {
			return node
		}
		return nil
	}
	n, ok := node.getChildren()[runes[0]]
	if !ok {
		return nil
	}
	var nrunes []rune
	if len(runes) > 1 {
		nrunes = runes[1:]
	} else {
		nrunes = runes[0:0]
	}
	return findLastNode(n, nrunes, value)
}
func (node trieNode) getParent() *trieNode {
	return node.parent
}

// Returns the children of this node.
func (node trieNode) getChildren() map[rune]*trieNode {
	return node.children
}

func searchSlice(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
