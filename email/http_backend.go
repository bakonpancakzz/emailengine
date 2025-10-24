package email

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func newHttpHandler(e *Engine) *http.ServeMux {
	v := validator.New()
	r := http.NewServeMux()
	r.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {

		// Sanity Checks
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		if r.ContentLength > e.IncomingMaxBytes {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		if !e.AuthHandler(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse Request Body
		var incoming []Email
		var incomingBody io.Reader
		var consumedBody = http.MaxBytesReader(w, r.Body, e.IncomingMaxBytes)

		switch strings.ToLower(r.Header.Get("Content-Encoding")) {
		case "gzip":
			gr, err := gzip.NewReader(consumedBody)
			if err != nil {
				// Errors only arise from invalid gzip headers so this is
				// definitely the users fault
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer gr.Close()
			incomingBody = gr

		case "": // No Compression
			incomingBody = consumedBody

		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		decoder := json.NewDecoder(incomingBody)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&incoming); err != nil {
			e.ErrorLogger(fmt.Errorf("error parsing body: %s", err))
			http.Error(w, "Invalid Form Body", http.StatusUnprocessableEntity)
			return
		}

		// Queue Incoming Emails
		for i := range incoming {
			if err := v.Struct(incoming[i]); err != nil {
				e.ErrorLogger(fmt.Errorf("validation failed for email at index %d: %s", i, err))
				http.Error(w, fmt.Sprintf("Validation Failed for Email at Index %d: %s\n", i, err), http.StatusBadRequest)
				continue
			}
			if ok := e.QueueEmail(&incoming[i]); !ok {
				http.Error(w, fmt.Sprintf("Email queue is full at index: %d\n", i), http.StatusInsufficientStorage)
				continue
			}
		}

		// Success!
		w.WriteHeader(http.StatusCreated)
	})
	return r
}
