package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// GetSession will return an error if there's no session or error while unmarshal data session
//
// @param dest interface{}
// @param req *http.Request
func GetSession(dest interface{}, req *http.Request) error {
	var (
		session = req.Header.Get("x-sess-token")
		err     error
	)

	if reflect.ValueOf(session).IsZero() {
		err = fmt.Errorf("session not found")
		return err
	}
	
	err = json.Unmarshal([]byte(session), dest)
	return err
}
