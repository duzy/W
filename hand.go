package W

import (
        "net/http"
        "fmt"
        "os"
)

type DealContext struct {
        RW http.ResponseWriter
        R *http.Request
        Data interface{}
}

func (ren *HTMLRen) Deal(name string, dealer func(*DealContext)) http.Handler {
        t := ren.T
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                dc := &DealContext{ w, r, nil }
                dealer(dc)
                ren.closure = dc
                if err := t.ExecuteTemplate(w, name, dc.Data); err != nil {
                        fmt.Fprintf(os.Stderr, "%v\n", err)
                }
                ren.closure = nil
        })
}

func Deal(name string, dealer func(*DealContext)) http.Handler {
        return DefaultRen.Deal(name, dealer)
}
