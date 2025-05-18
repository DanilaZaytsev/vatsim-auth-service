package handler

import (
	"fmt"
	"net/http"
)

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать в админскую панель!")
}
