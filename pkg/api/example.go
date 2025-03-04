package api

import (
	"fmt"
	"net/http"
)

func ExampleHamdler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello haiiii")
}
