package simple_mitm

import (
	"crypto/tls"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"io"
	"net"
	"net/http"
	"time"
)

type NoNextTLSProto map[string]func(*http.Server, *tls.Conn, http.Handler)

func Handle(w http.ResponseWriter, r *http.Request) {
	for range only.Once {
		dest, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("cannot dial: %s", err),
				http.StatusServiceUnavailable,
			)
			break
		}
		w.WriteHeader(http.StatusOK)
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w,
				"Hijacking not supported",
				http.StatusInternalServerError,
			)
			break
		}
		client, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w,
				fmt.Sprintf("cannot hijack: %s", err),
				http.StatusInternalServerError,
			)
			break
		}

		go transfer(w, dest, client)
		go transfer(w, client, dest)
	}
}

//goland:noinspection GoUnusedParameter
func transfer(w http.ResponseWriter, dst io.WriteCloser, src io.ReadCloser) {
	defer closeConn(dst)
	defer closeConn(src)
	for range only.Once {
		_, err := io.Copy(dst, src)
		if err == nil {
			break
		}
		if isClosedNetConn(err) {
			break
		}
		//errorChannel <- errors.New(
		//	fmt.Sprintf("cannot copy from source to destination: %s\n", err))

	}
}

func isClosedNetConn(err error) bool {
	is := false
	for range only.Once {
		var ok bool
		if err, ok = err.(*net.OpError); !ok {
			break
		}
		if err.(*net.OpError).Op != "read" {
			break
		}
		if err.(*net.OpError).Net != "tcp" {
			break
		}
		if err.(*net.OpError).Err.Error() != "use of closed network connection" {
			break
		}
		is = true
	}
	return is
}
func closeConn(c io.Closer) {
	_ = c.Close()
}
