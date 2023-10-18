package healthz

import "net/http"

func Example() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckStaticMD(map[string]string{
		"key1": "value1",
	}))
	checker.AddCheck("app2", NewCheckStaticMD(map[string]string{
		"key2": "a different value",
	}))

	http.ListenAndServe(":8080", checker)
}

func Example_gorilla() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckStaticMD(map[string]string{
		"key1": "value1",
	}))
	checker.AddCheck("app2", NewCheckStaticMD(map[string]string{
		"key2": "a different value",
	}))

	http.Handle("/healthz", checker)

	http.ListenAndServe(":8080", nil)
}
