package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type signUpDTO struct {
	Email    string `json:"email"`
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dto signUpDTO
	err := decoder.Decode(&dto)
	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}
	fmt.Printf("email: %s\n", dto.Email)
	fmt.Printf("userId: %s\n", dto.UserId)
	fmt.Printf("password: %s\n", dto.Password)
	fmt.Fprintln(w, "OK")
}
