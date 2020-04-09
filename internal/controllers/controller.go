package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.unanet.io/devops/eve/internal/common"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// AppErr is syntax sugar to return an Err
func AppErr(err error, message string) (int, interface{}, error) {
	return 0, nil, err
}

// AppResponse is a syntax sugar to return a Response
func AppResponse(code int, response interface{}) (int, interface{}, error) {
	return code, response, nil
}

// Sugar to de-dupe usage below
// be sure an return whenever you write
// or else multiple writes can go out
func writeResponse(w http.ResponseWriter, status int, response []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(response)
	return
}

// Handler is a wrapper for routes to pick up errors and return them appropriately
type Handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, payload, err := fn(w, r)
	if err == nil {
		if payload != nil {
			response, _ := json.Marshal(payload)
			writeResponse(w, status, response)
			return
		}
		w.WriteHeader(status)
		return
	}

	var restError *common.RestError
	if errors.As(err, &restError) {
		response, _ := json.Marshal(restError)
		writeResponse(w, restError.Code, response)
		return
	}

	response, _ := json.Marshal(common.RestError{
		Code:    500,
		Message: fmt.Sprintf("type: %s, %+v\n", reflect.TypeOf(err), err),
	})
	writeResponse(w, http.StatusInternalServerError, response)
	return

}

// Controller is a base Controller
type Controller struct {
}

// Vars gets the url variables
func (c *Controller) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Query gets the query parameter values
func (c *Controller) Query(r *http.Request) map[string][]string {
	return r.URL.Query()
}

// ParseBody of incoming request
func (c *Controller) ParseBody(res http.ResponseWriter, r *http.Request, model interface{}) error {
	// json decode the payload - obviously this could be abstracted
	// to handle many content types
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		return &common.RestError{
			Code:    400,
			Message: fmt.Sprintf("Invalid Post Body: %s", err),
		}
	}

	defer r.Body.Close()
	return nil
}
