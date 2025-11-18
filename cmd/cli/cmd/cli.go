package cmd

import (
	"fmt"

	"github.com/google/uuid"
)

func createUUID() uuid.UUID {
	id := uuid.New()
	fmt.Print(id)
	return id
}
