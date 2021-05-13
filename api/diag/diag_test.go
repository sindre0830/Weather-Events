package diag

import (
	"main/dict"

	"net/http"
	"testing"
)

func TestRequestHandler(t *testing.T) {

	http.HandleFunc(dict.DIAG_PATH, MethodHandler)

}
