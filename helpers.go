package vkstream

import (
	"errors"
	"fmt"
)

func newVkStreamingError(error vkStreamingError) error {
	return errors.New(fmt.Sprintf("%d: %s", error.ErrorCode, error.Message))
}

func newVkError(error vkError) error {
	return errors.New(fmt.Sprintf("%d: %s", error.ErrorCode, error.Message))
}
