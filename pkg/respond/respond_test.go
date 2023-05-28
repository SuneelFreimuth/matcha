package respond

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)


func TestProblemToJSON(t *testing.T) {
	t.Run("empty problem -> empty object", func(t *testing.T) {
		p := Problem{}
		bs, err := json.Marshal(p)
		if err != nil {
			t.Error("error while encoding Problem as JSON:", err)
		}
		if !reflect.DeepEqual(bs, []byte("{}")) {
			t.Errorf("empty problem did not encode as empty object, got %s", bs)
		}
	})

	t.Run("complete problem -> complete object", func(t *testing.T) {
		p := Problem{
			Type: "https://example.io/problems/out-of-credit",
			Title: "You do not have enough credit.",
			Detail: "Your current balance is 30, but that costs 50.",
			Status: 418,
			Instance: "/account/12345/msgs/abc",
			Extensions: map[string]any{
				"balance": 30,
				"accounts": []string{
					"/account/12345",
					"/account/67890",
				},
			},
		}
		bs, err := json.Marshal(p)
		if err != nil {
			t.Error("error while encoding Problem as JSON:", err)
		}
		expected := `{"type":"https://example.io/problems/out-of-credit","title":"You do not have enough credit.",` +
			`"status":418,"detail":"Your current balance is 30, but that costs 50.",` + 
			`"instance":"/account/12345/msgs/abc",` +
			`"balance":30,"accounts":["/account/12345","/account/67890"]}`
		if !reflect.DeepEqual(bs, []byte(expected)) {
			fmt.Println("incorrect encoding of Problem:")
			fmt.Printf("  Expected: %s\n", expected)
			fmt.Printf("  Received: %s\n", bs)
			t.Fail()
		}
	})
}

func TestProblemFromJSON(t *testing.T) {
	t.Run("empty JSON object -> empty Problem", func (t *testing.T)  {
		var p Problem
		err := json.Unmarshal([]byte("{}"), &p)
		if err != nil {
			t.Errorf("error while unmarshaling JSON, %v", err)
		}
		if !reflect.DeepEqual(p, Problem{}) {
			t.Errorf("problem should be empty, got %#v", p)
		}
	})
	t.Run("complete JSON object -> complete Problem", func (t *testing.T)  {
		var p Problem
		err := json.Unmarshal([]byte(`{
			"type": "https://example.io/problems/out-of-credit",
			"title": "You do not have enough credit.",
			"status": 418,
			"detail": "Your current balance is 30, but that costs 50.",
			"instance": "/account/12345/msgs/abc",
			"balance": 30,
			"accounts": [
				"/account/12345",
				"/account/67890"
			]
		}`), &p)
		if err != nil {
			t.Errorf("error while unmarshaling JSON, %v", err)
		}

		expected := Problem{
			Type: "https://example.io/problems/out-of-credit",
			Title: "You do not have enough credit.",
			Detail: "Your current balance is 30, but that costs 50.",
			Status: 418,
			Instance: "/account/12345/msgs/abc",
			Extensions: map[string]any{
				"balance": float64(30),
				"accounts": []any{
					"/account/12345",
					"/account/67890",
				},
			},
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf(
				"problem did not deserialize successfully.\nReceived: %#v\nExpected: %#v\n",
				p,
				expected,
			)
		}
	})
}