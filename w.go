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
        dc *DealContext
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

func (ren *HTMLRen) Delims(left, right string) *HTMLRen {
        ren.T.Delims(left, right)
        return ren
}

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

func (ren *HTMLRen) LoadString(name string) (s string, err error) {
        if t := ren.T.Lookup(name); t != nil {
                w := new(bytes.Buffer)
                if err = t.Execute(w, nil); err == nil {
                        s = w.String()
                }
        }
        return
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

func (ren *HTMLRen) MustLoadString(name string) (res string) {
        if s, err := ren.LoadString(name); err != nil {
                panic(err)
        } else {
                res = s
        }
        return
}

func (ren *HTMLRen) CanExpand(name string) bool {
        return ren.T.Lookup(name) != nil
}

func Delims(left, right string) *HTMLRen {
        return DefaultRen.Delims(left, right)
}

func MustGlob(patterns ...string) {
        DefaultRen.MustGlob(patterns...)
}

func MustExpand(w io.Writer, name string) {
        DefaultRen.MustExpand(w, name)
}

func MustLoadString(name string) string {
        return DefaultRen.MustLoadString(name)
}

func CanExpand(name string) bool {
        return DefaultRen.CanExpand(name)
}

func closureFunc(ren *HTMLRen) interface{} {
        return func(name string, args ...string) template.HTML {
                // TODO: buffer pool for performance
                w, cc := new(bytes.Buffer), newClosureContext(args...)
                cc.Set("Context", ren.dc.closure)
                if err := ren.expand(w, name, cc); err != nil {
                        panic(err) // FIXME: error handling?
                }
                return template.HTML(w.String())
        }
}

type ClosureContext map[string]interface{}

func (cc ClosureContext) Set(name string, data interface{}) {
        cc[name] = data
}

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
