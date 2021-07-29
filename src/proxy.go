package simple_mitm

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"net/http"
)

type Proxy struct {
	*http.Server
}

// NewProxy returns a new instance of Proxy
func NewProxy() *Proxy {
	return &Proxy{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%s", ProxyPort),
			Handler:      handleHTTPSFunc(),
			TLSNextProto: make(NoNextTLSProto), // Disable HTTP/2.
		},
	}
}

func handleHTTPSFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for range only.Once {
			if r.Method != http.MethodConnect {
				HandleHTTPS(w, r)
				break
			}
			r.URL.Host = APIServerAndPort
			r.Host = APIServerAndPort
			Handle(w, r)
		}
	})
}
