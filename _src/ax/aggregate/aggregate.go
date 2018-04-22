package aggregate

import (
	"github.com/jmalloc/ax/src/ax/ident"
)

// ID uniquely identifies an aggregate.
type ID struct{ ident.ID }

// TODO
type Aggregate interface{}
