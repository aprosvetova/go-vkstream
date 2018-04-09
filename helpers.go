package vkstream

import (
	"fmt"
)

func newVkStreamingError(error vkStreamingError) error {
	return fmt.Errorf("%d: %s", error.ErrorCode, error.Message)
}

func newVkError(error vkError) error {
	return fmt.Errorf("%d: %s", error.ErrorCode, error.Message)
}
