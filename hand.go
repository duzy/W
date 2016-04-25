package W

import (
        "net/http"
        "fmt"
        "os"
)

type DealContext struct {
        RW http.ResponseWriter
        R *http.Request
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
        t := ren.T
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                ren.dc = &DealContext{ w, r, nil }
                dealer(ren.dc)
                if err := t.ExecuteTemplate(w, name, ren.dc.closure); err != nil {
                        fmt.Fprintf(os.Stderr, "%v\n", err)
                }
        })
}

func Deal(name string, dealer func(*DealContext)) http.Handler {
        return DefaultRen.Deal(name, dealer)
}
