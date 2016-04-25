package W

import (
        "html/template"
        "strings"
        "bytes"
        "io"
)

var DefaultRen = NewHTMLRen("DEFAULT")

type HTMLRen struct {
        T *template.Template
        closure *DealContext
}

func NewHTMLRen(name string) (ren *HTMLRen) {
        ren = &HTMLRen{
                T: template.New(name),
        }
        ren.T.Funcs(initialFuncMap(ren))
        return
}

func initialFuncMap(ren *HTMLRen) template.FuncMap {
        return template.FuncMap{
                "closure": closureFunc(ren),
        }
}

/*
func NewHTMLRenWithTemplate(t *template.Template) *HTMLRen {
        return &HTMLRen{
                T: t,
        }
} */

func (ren *HTMLRen) Glob(pattern string) (*HTMLRen, error) {
        _, err := ren.T.ParseGlob(pattern)
        return ren, err
}

func (ren *HTMLRen) Expand(w io.Writer, name string) error {
        return ren.expand(w, name, nil) // FIXME: parssing some default data
}

func (ren *HTMLRen) expand(w io.Writer, name string, data interface{}) error {
        if len(name) == 0 {
                return ren.T.Execute(w, data)
        }
        return ren.T.ExecuteTemplate(w, name, data)
}

func (ren *HTMLRen) MustGlob(patterns ...string) *HTMLRen {
        for _, s := range patterns {
                if _, err := ren.Glob(s); err != nil {
                        panic(err)
                }
        }
        return ren
}

func (ren *HTMLRen) MustExpand(w io.Writer, name string) *HTMLRen {
        if err := ren.Expand(w, name); err != nil {
                panic(err)
        }
        return ren
}

func MustGlob(patterns ...string) {
        DefaultRen.MustGlob(patterns...)
}

func MustExpand(w io.Writer, name string) {
        DefaultRen.MustExpand(w, name)
}

func closureFunc(ren *HTMLRen) interface{} {
        return func(name string, args ...string) template.HTML {
                // TODO: buffer pool for performance
                w, cc := new(bytes.Buffer), newClosureContext(args...)
                if err := ren.expand(w, name, cc); err != nil {
                        panic(err) // FIXME: error handling?
                }
                return template.HTML(w.String())
        }
}

// TODO: support {{ .Context.XXX }}
type ClosureContext map[string]interface{}

// TODO: a better and stronger closure syntax
func newClosureContext(args ...string) ClosureContext {
        cc := make(ClosureContext)
        for _, a := range args {
                if i := strings.Index(a, "="); 0 < i {
                        name := strings.TrimSpace(a[0:i])
                        value := strings.TrimSpace(a[i+1:])
                        if name != "" {
                                cc[name] = value
                        }
                }
        }
        return cc
}
