package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Problem interface {
	GetStatus() int
}

// The minimum implementation of a problem details object conforming to RFC 7807.
//
// To define a problem details object with extensions, embed BaseProblem:
// 
//  type purchaseProblem struct {
//      respond.BaseProblem
//      Balance  int      `json:"balance,omitempty"`
//      Accounts []string `json:"accounts,omitempty"`
//  }
type BaseProblem struct {
	Type     string `json:"type,omitempty" xml:"type,omitempty"`
	Title    string `json:"title,omitempty" xml:"title,omitempty"`
	Status   int    `json:"status,omitempty" xml:"status,omitempty"`
	Detail   string `json:"detail,omitempty" xml:"detail,omitempty"`
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`
}

func (bp BaseProblem) GetStatus() int {
	return bp.Status
}

// Write an HTTP response with the JSON encoding of the problem as the response body.
//
// The status code will be taken from the problem. If absent, it will be 400 Bad Request.
func WithProblem(w http.ResponseWriter, p Problem) error {
	var bs []byte
	var err error
	w.Header().Set("Content-Type", "application/problem+json")
	bs, err = json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to encode Problem as JSON: %v", err)
	}

	if s := p.GetStatus(); s != 0 {
		w.WriteHeader(s)
	} else {
		w.WriteHeader(400)
	}

	_, err = w.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
