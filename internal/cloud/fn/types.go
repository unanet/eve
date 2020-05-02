package fn

import (
	"github.com/dghubble/sling"
)

type Argument func(*sling.Sling) *sling.Sling

func Azure(functionUrl, code string) Argument {
	return func(s *sling.Sling) *sling.Sling {
		return s.Base(functionUrl).QueryStruct(struct {
			code string
		}{code: code})
	}
}

func MapBody(body map[string]interface{}) Argument {
	return func(s *sling.Sling) *sling.Sling {
		return s.BodyJSON(body)
	}
}
