package srv

import (
	"google.golang.org/grpc/naming"
)

// ResolverFunc implements the naming.Resolver interface
type ResolverFunc func(target string) (naming.Watcher, error)

// Resolve calls the ResolverFunc
func (f ResolverFunc) Resolve(target string) (naming.Watcher, error) {
	return f(target)
}

func resolver(target string) (naming.Watcher, error) {
	return NewWatcher(target), nil
}

// Resolver implements the ResolverFunc interface
var Resolver = ResolverFunc(resolver)
