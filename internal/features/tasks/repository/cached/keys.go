package tasks_cached_repository

import (
	"strconv"

	"github.com/google/uuid"
)

func taskKey(id uuid.UUID) string {
	return "task:" + id.String()
}

func tasksListKey(userID *uuid.UUID) string {
	if userID == nil {
		return "tasks:all"
	}

	return "tasks:" + userID.String()
}

func tasksListField(limit, offset *int) string {
	format := func(value *int) string {
		if value == nil {
			return "nil"
		}

		return strconv.Itoa(*value)
	}

	return format(limit) + ":" + format(offset)
}
