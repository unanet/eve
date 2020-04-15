package fn

import (
	"fmt"

	"github.com/dghubble/sling"
)

type Argument func(*sling.Sling) *sling.Sling

func Azure(functionApp, function, code string) Argument {
	return func(s *sling.Sling) *sling.Sling {
		return s.Base(fmt.Sprintf("https://%s.azurewebsites.net/api/%s", functionApp, function)).QueryStruct(struct {
			code string
		}{code: code})
	}
}

func MapBody(body map[string]interface{}) Argument {
	return func(s *sling.Sling) *sling.Sling {
		return s.BodyJSON(body)
	}
}
