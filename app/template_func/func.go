package template_func

import (
	"html/template"
	"strings"
)

type TemplateFunc struct{}

func NewTemplateFunc() *TemplateFunc {
	return &TemplateFunc{}
}

func (t *TemplateFunc) FuncMap() template.FuncMap {
	return template.FuncMap{
		"unsafe": func(str string) template.HTML {
			return template.HTML(str)
		},
		"upper": func(str string) string {
			return strings.ToUpper(str)
		},
		"lower": func(str string) string {
			return strings.ToLower(str)
		},
	}
}
