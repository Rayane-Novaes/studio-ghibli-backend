package app

import (
	"errors"
	"io"
	"log"
	"net/http"
)

func Run() {

	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func echo (w http.ResponseWriter, r *http.Request){
	var body []byte
		buffer := make([]byte, 4)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		for {
			num, err := r.Body.Read(buffer)

			if num > 0 {
				body = append(body, buffer[:num]...)
			}

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	
}

