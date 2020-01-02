package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/justinpage/go-web-calculator/evaluator"
)

func main() {
	port := flag.String("port", "8080", "port number")
	flag.Parse()

	http.HandleFunc("/", render)
	http.HandleFunc("/calculate", calculate)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func render(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func calculate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	expr, err := parseAndCheck(r.Form.Get("expr"))
	if err != nil {
		http.Error(w, "bad expr: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%.6g", expr.Eval(evaluator.Env{}))
}

func parseAndCheck(s string) (evaluator.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	expr, err := evaluator.Parse(s)
	if err != nil {
		return nil, err
	}
	vars := make(map[evaluator.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}
	return expr, nil
}
