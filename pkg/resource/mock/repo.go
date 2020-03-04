package mock

import (
	"context"
	"fmt"
	"strings"

	"github.com/hanjunlee/argocui/pkg/resource"
)

// Animal is the name of animal.
type Animal string

var (
	zoo = []Animal{"alligator", "ant", "bear", "bee", "camel", "cat", "cheetah"}
)

// Repo is the mocking repository.
type Repo struct{}

// Get return animal.
func (r *Repo) Get(key string) interface{} {
	for _, a := range zoo {
		if string(a) == key {
			return a
		}
	}

	return ""
}

// Search return animal which is matched with the pattern.
func (r *Repo) Search(pattern string) []interface{} {
	animals := make([]interface{}, 0)
	for _, a := range zoo {
		if i := strings.Index(string(a), pattern); i != -1 {
			animals = append(animals, a)
		}
	}
	return animals
}

// Delete delete the animal.
func (r *Repo) Delete(key string) error {
	for idx, a := range zoo {
		if string(a) != key {
			continue
		}

		if idx == len(zoo) - 1 {
			zoo = zoo[0:idx]
		} else {
			zoo = append(zoo[:idx], zoo[idx+1:len(zoo)]...)
		}
		return nil
	}
	
	return fmt.Errorf("it doesn't exist: %s", key)
}

// Logs is not implemented.
func (r *Repo) Logs(ctx context.Context, key string) (ch <-chan resource.Log, err error) {
	return nil, resource.ErrNotImplement
}
