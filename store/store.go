package store

import (
	_ "github.com/lib/pq"
)

type Dependency struct {
	QuestionStore QuestionStore
}
