package response

import (
	"context"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	Json().Success(context.TODO(), "")
	return
}
