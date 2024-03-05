package wildlifenl

import (
	"net/http"
)

type Me struct {
	Controller
}

func newMe() *Me {
	controller := new(Me)
	controller.getFunc = controller.get
	return controller
}

func (c *Me) get(r *http.Request) ([]byte, httpStatus, error) {
	data := []byte("{\"func\":\"me.get\"}")
	return data, http.StatusOK, nil
}
