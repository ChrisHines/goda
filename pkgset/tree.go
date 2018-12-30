package pkgset

import (
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Tree struct {
	Path    string
	Package *packages.Package

	Child map[string]*Tree

	Parent   *Tree
	Children []*Tree
}

func NewTree(parent *Tree, path string) *Tree {
	return &Tree{
		Path:   path,
		Child:  map[string]*Tree{},
		Parent: parent,
	}
}

func (tree *Tree) Add(pkg *packages.Package) {
	tree.Insert([]string{}, strings.Split(pkg.PkgPath, "/"), pkg)
}

func (tree *Tree) Insert(prefix, suffix []string, pkg *packages.Package) {
	if len(suffix) == 0 {
		tree.Package = pkg
		return
	}

	childPrefix := append(prefix, suffix[0])
	child, hasChild := tree.Child[suffix[0]]
	if !hasChild {
		child = NewTree(tree, path.Join(childPrefix...))
		tree.Child[suffix[0]] = child
		tree.Children = append(tree.Children, child)
	}

	child.Insert(childPrefix, suffix[1:], pkg)
}

func (tree *Tree) Walk(fn func(tree *Tree)) {
	fn(tree)
	for _, child := range tree.Children {
		child.Walk(fn)
	}
}
