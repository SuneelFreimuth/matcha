package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A problem details object conforming to RFC 7807 §3.1.
// https://datatracker.ietf.org/doc/html/rfc7807#section-3.1
// 
// If Type, Title, Detail, or Instance is an empty string, it will not be included
// in the JSON/XML encoding of Problem.
// 
//
// TODO: Examples.
type Problem struct {
	Type string
	Title string
	Status int
	Detail string
	Instance string
	Extensions map[string]any
}

func (p Problem) MarshalJSON() ([]byte, error) {
	buf := []byte{ '{' }
	if p.Type != "" {
		buf = append(buf, fmt.Sprintf(`"type":"%s",`, p.Type)...)
	}
	if p.Title != "" {
		buf = append(buf, fmt.Sprintf(`"title":"%s",`, p.Title)...)
	}
	if p.Status != 0 {
		buf = append(buf, fmt.Sprintf(`"status":%d,`, p.Status)...)
	}
	if p.Detail != "" {
		buf = append(buf, fmt.Sprintf(`"detail":"%s",`, p.Detail)...)
	}
	if p.Instance != "" {
		buf = append(buf, fmt.Sprintf(`"instance":"%s",`, p.Instance)...)
	}
	for key, value := range p.Extensions {
		value, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		buf = append(buf, fmt.Sprintf(`"%s":%s,`, key, value)...)
	}
	if buf[len(buf)-1] == ',' {
		buf[len(buf)-1] = '}'
	} else {
		buf = append(buf, '}')
	}
	return buf, nil
}

func WithProblem(w http.ResponseWriter, p Problem) error {
	return json.NewEncoder(w).Encode(p) 
}