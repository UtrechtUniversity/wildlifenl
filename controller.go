package wildlifenl

import (
	"log"
	"net/http"
)

type httpStatus int

type httpHandler func(r *http.Request) ([]byte, httpStatus, error)

type Controller struct {
	getFunc    httpHandler
	postFunc   httpHandler
	putFunc    httpHandler
	deleteFunc httpHandler
}

func (c *Controller) handle(w http.ResponseWriter, r *http.Request) {
	if !app.authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var data []byte
	var status httpStatus
	var err error
	switch r.Method {
	case http.MethodGet:
		if c.getFunc == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		data, status, err = c.getFunc(r)
	case http.MethodPost:
		if c.postFunc == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		data, status, err = c.postFunc(r)
	case http.MethodPut:
		if c.putFunc == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		data, status, err = c.putFunc(r)
	case http.MethodDelete:
		if c.deleteFunc == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		data, status, err = c.deleteFunc(r)
	}
	if err != nil {
		log.Println("ERROR handling", r.Method, ":", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(int(status))
	writeResponseJSON(w, data)
}
