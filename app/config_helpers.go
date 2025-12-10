package app

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type YamlComments struct{}

func MyParser() *YamlComments {
	return &YamlComments{}
}

func (p *YamlComments) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if out["general"] == nil {
		return out, fmt.Errorf("invalid config file: general key not found")
	}

	return out, nil
}

func (p *YamlComments) Marshal(o map[string]interface{}) ([]byte, error) {
	comments := []*yaml.Comment{{Texts: []string{"This is a file-level comment"}}}
	cm := yaml.CommentMap{
		"x": comments,
	}
	data, err := yaml.MarshalWithOptions(o, yaml.WithComment(cm))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal with comment: %w", err)
	}
	return data, nil
}
