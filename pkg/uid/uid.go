package uid

import (
	"github.com/google/uuid"
)

type Generator struct {
}

func NewGenerator() Generator {
	return Generator{}
}

func (generator Generator) NewUUID() string {
	return uuid.New().String()
}
