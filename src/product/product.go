package product

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<h1>"+vars["productId"]+"</h1>")
	w.WriteHeader(http.StatusOK)
}
