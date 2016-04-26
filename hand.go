package W

import (
        "net/http"
        "fmt"
        "os"
)

type DealContext struct {
        RW http.ResponseWriter
        R *http.Request
        Name string
        closure ClosureContext
}

// FIXME: thread safety?
func (dc *DealContext) Set(name string, data interface{}) {
        if dc.closure == nil {
                dc.closure = make(ClosureContext)
        }
        dc.closure.Set(name, data)
}

func (ren *HTMLRen) Deal(name string, dealer func(*DealContext)) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                ren.dc = &DealContext{ w, r, name, nil }
                dealer(ren.dc); name = ren.dc.Name
                if err := ren.T.ExecuteTemplate(w, name, ren.dc.closure); err != nil {
                        fmt.Fprintf(os.Stderr, "%v\n", err)
                }
                ren.dc = nil
        })
}

func Deal(name string, dealer func(*DealContext)) http.Handler {
        return DefaultRen.Deal(name, dealer)
}
