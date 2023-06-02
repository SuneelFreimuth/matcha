package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)


type purchaseProblem struct {
	BaseProblem
	Balance int `json:"balance"`
	Accounts []string `json:"accounts"`
}


func TestProblemToJSON(t *testing.T) {
	t.Run("empty problem -> only extensions", func(t *testing.T) {
		p := purchaseProblem{}
		bs, err := json.Marshal(p)
		if err != nil {
			t.Error("error while encoding Problem as JSON:", err)
		}
		if !reflect.DeepEqual(bs, []byte(`{"balance":0,"accounts":null}`)) {
			t.Errorf("empty problem did not only encode the extensions, got %s", bs)
		}
	})

	t.Run("complete problem -> complete object", func(t *testing.T) {
		p := purchaseProblem{
			BaseProblem: BaseProblem{
				Type: "https://example.io/problems/out-of-credit",
				Title: "You do not have enough credit.",
				Detail: "Your current balance is 30, but that costs 50.",
				Status: 418,
				Instance: "/account/12345/msgs/abc",
			},
			Balance: 30,
			Accounts: []string{
				"/account/12345",
				"/account/67890",
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
		var p purchaseProblem
		err := json.Unmarshal([]byte("{}"), &p)
		if err != nil {
			t.Errorf("error while unmarshaling JSON, %v", err)
		}
		if !reflect.DeepEqual(p, purchaseProblem{}) {
			t.Errorf("problem should be empty, got %#v", p)
		}
	})
	t.Run("complete JSON object -> complete Problem", func (t *testing.T)  {
		var p purchaseProblem
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

		expected := purchaseProblem{
			BaseProblem: BaseProblem{
				Type: "https://example.io/problems/out-of-credit",
				Title: "You do not have enough credit.",
				Detail: "Your current balance is 30, but that costs 50.",
				Status: 418,
				Instance: "/account/12345/msgs/abc",
			},
			Balance: 30,
			Accounts: []string{
				"/account/12345",
				"/account/67890",
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

type testResponseWriter struct {
	*strings.Builder
	header http.Header
	Status int
}

func (w testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
}

func TestWithProblem(t *testing.T) {
	w := testResponseWriter{
		Builder: &strings.Builder{},
		header: make(http.Header, 0),
	}
	p := purchaseProblem{
		BaseProblem: BaseProblem{
			Type: "https://example.io/problems/out-of-credit",
			Title: "You do not have enough credit.",
			Detail: "Your current balance is 30, but that costs 50.",
			Status: 418,
			Instance: "/account/12345/msgs/abc",
		},
		Balance: 30,
		Accounts: []string{
			"/account/12345",
			"/account/67890",
		},
	}
	WithProblem(&w, p)

	if w.Status != p.Status {
		t.Errorf("response status code did not match expected. received %d, expected %d", w.Status, p.Status)
	}

	expected := `{"type":"https://example.io/problems/out-of-credit","title":"You do not have enough credit.",` +
			`"status":418,"detail":"Your current balance is 30, but that costs 50.",` + 
			`"instance":"/account/12345/msgs/abc",` +
			`"balance":30,"accounts":["/account/12345","/account/67890"]}`
	if result := w.String(); result != expected {
		t.Errorf(
			"problem did not serialize successfully.\nReceived: %#v\nExpected: %#v\n",
			result,
			expected,
		)
	}
}