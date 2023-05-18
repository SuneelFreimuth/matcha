package respond

import (
	"fmt"
	"encoding/json"
	"strings"
	"testing"
)

func TestProblemToJSON(t *testing.T) {
	t.Run("empty problem -> empty object", func(t *testing.T) {
		p := Problem{}
		var result strings.Builder
		err := json.NewEncoder(&result).Encode(p)
		if err != nil {
			t.Error("error while encoding Problem as JSON:", err)
		}
		s := result.String()
		fmt.Printf("%x\n", s)
		if s != "{}" {
			t.Error("empty problem did not encode as empty object, got", s)
		}
	})
	t.Run("complete problem -> complete object", func(t *testing.T) {
		p := Problem{
			Type: "https://example.io/problems/out-of-credit",
			Title: "You do not have enough credit.",
			Detail: "Your current balance is 30, but that costs 50.",
			Instance: "/account/12345/msgs/abc",
			Extensions: map[string]any{
				"balance": 30,
				"accounts": []string{
					"/account/12345",
					"/account/67890",
				},
			},
		}
	})
}