package healthz

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Example() {
	checker := NewChecker()
	checker.AddCheck("app", NewStaticMD(map[string]string{
		"key1": "value1",
	}))
	checker.AddCheck("app2", NewStaticMD(map[string]string{
		"key2": "a different value",
	}))

	http.ListenAndServe(":8080", checker)
}

func Example_gorilla() {
	checker := NewChecker()
	checker.AddCheck("app", NewStaticMD(map[string]string{
		"key1": "value1",
	}))
	checker.AddCheck("app2", NewStaticMD(map[string]string{
		"key2": "a different value",
	}))

	r := mux.NewRouter()
	r.PathPrefix("/healthz").Handler(checker)

	http.ListenAndServe(":8080", r)
}
