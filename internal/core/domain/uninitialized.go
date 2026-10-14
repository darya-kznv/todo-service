package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	UninitializedID        = uuid.Nil
	UninitializedVersion   = -1
	UninitializedCreatedAt = time.Time{}
)
