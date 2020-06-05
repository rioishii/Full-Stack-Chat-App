package indexes

import (
	"reflect"
	"testing"
)

func TestTrieAdd(t *testing.T) {
	cases := []struct {
		name         string
		keys         []string
		values       []int64
		expectedSize int
	}{
		{
			"Inserting distinct sets",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			6,
		},
		{
			"Inserting duplicate",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 1, 4, 5},
			5,
		},
	}

	for _, c := range cases {
		trie := NewTrie()
		len := len(c.keys)

		for i := 0; i < len; i++ {
			trie.Add(c.keys[i], c.values[i])
		}

		if trie.Len() != c.expectedSize {
			t.Errorf("case %s: expected trie size: %v (got %v)", c.name, c.expectedSize, trie.Len())
		}
	}
}

func TestTrieFind(t *testing.T) {
	trie := NewTrie()
	trie.Add("go", 1)
	trie.Add("git", 2)
	trie.Add("gob", 3)
	trie.Add("go", 4)
	trie.Add("goal", 5)
	trie.Add("foo", 1)

	cases := []struct {
		name           string
		prefix         string
		max            int
		expectedValues []int64
	}{
		{
			"case 1",
			"go",
			10,
			[]int64{1, 4, 3, 5},
		},
		{
			"case 2",
			"f",
			1,
			[]int64{1},
		},
		{
			"case 3",
			"x",
			10,
			nil,
		},
		{
			"case 4",
			"go",
			2,
			[]int64{1, 4},
		},
	}

	for _, c := range cases {

		values := trie.Find(c.prefix, c.max)

		if !reflect.DeepEqual(values, c.expectedValues) {
			t.Errorf("case %s: expected values: %v (got %v)", c.name, c.expectedValues, values)
		}
	}
}

func TestTrieRemove(t *testing.T) {
	trie := NewTrie()
	trie.Add("go", 1)
	trie.Add("git", 2)
	trie.Add("gob", 3)
	trie.Add("go", 4)
	trie.Add("goal", 5)
	trie.Add("foo", 1)

	cases := []struct {
		name           string
		delete         string
		id             int64
		find           string
		expectedValues []int64
		expectedSize   int
	}{
		{
			"case 1",
			"foo",
			1,
			"foo",
			nil,
			5,
		},
		{
			"case 2",
			"gob",
			3,
			"go",
			[]int64{1, 4, 5},
			4,
		},
		{
			"case 3",
			"unknown",
			1,
			"go",
			[]int64{1, 4, 5},
			4,
		},
		{
			"case 4",
			"go",
			4,
			"go",
			[]int64{1, 5},
			3,
		},
		{
			"case 5",
			"goal",
			20,
			"go",
			[]int64{1, 5},
			3,
		},
	}

	for _, c := range cases {
		trie.Remove(c.delete, c.id)
		values := trie.Find(c.find, 10)

		if !reflect.DeepEqual(values, c.expectedValues) {
			t.Errorf("case %s: expected values: %v (got %v)", c.name, c.expectedValues, values)
		}

		if trie.Len() != c.expectedSize {
			t.Errorf("case %s: expected trie size: %v (got %v)", c.name, c.expectedSize, trie.Len())
		}
	}

}
