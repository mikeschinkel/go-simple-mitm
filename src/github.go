package simple_mitm

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type API struct {
	*http.Server
}

// NewAPI returns a new instance of API
func NewAPI() *API {
	api := &API{}
	api.Server = &http.Server{
		Addr:         fmt.Sprintf(":%s", APIServerPort),
		TLSNextProto: make(NoNextTLSProto), // Disable HTTP/2.
		Handler:      api,
	}
	return api
}

//goland:noinspection GoUnusedParameter
func authGET(w http.ResponseWriter, r *http.Request) (resp *http.Response, err error) {
	for range only.Once {

		token := getAccessToken(r)

		if token == "" {
			err = fmt.Errorf("access token not found in HTTP '%s' header",
				HTTPAuthHeader)
			break
		}

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

		client := oauth2.NewClient(r.Context(), ts)
		r.TLS = nil       // Remove prior TLS object w/HandshakeComplete==true
		r.RequestURI = "" // Client.Do() will fail w/o this
		resp, err = client.Do(r)

	}
	return resp, err
}

func getAccessToken(r *http.Request) (token string) {
	for range only.Once {
		auth := strings.Split(r.Header.Get(HTTPAuthHeader), " ")
		if len(auth) != 2 {
			break
		}
		if auth[0] != HTTPAuthBearerToken {
			break
		}
		token = auth[1]
	}
	return token
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleHTTPS(w, r)
}

func HandleHTTPS(w http.ResponseWriter, r *http.Request) {
	for range only.Once {
		r.URL.Scheme = HTTPSScheme
		r.Host = GitHubAPIAndPort
		r.URL.Host = GitHubAPIAndPort
		resp, err := authGET(w, r)
		log.Printf("%v", resp)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("failed attempting authGET: %s", err),
				http.StatusServiceUnavailable,
			)
			break
		}
		// Write response to disk to allow for loading it instead of
		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		bytes, err := maybeSaveResponse(resp)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("cannot save response body: %s", err),
				http.StatusServiceUnavailable,
			)
			break
		}
		_, err = w.Write(bytes)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("cannot copy response body: %s", err),
				http.StatusServiceUnavailable,
			)
			break
		}
		_ = resp.Body.Close()
	}
}

type persistableResponse struct {
	StatusCode int         `json:"status"`
	Body       string      `json:"body"`
	Headers    http.Header `json:"headers"`
}

func maybeSaveResponse(resp *http.Response) (bytes []byte, err error) {
	for range only.Once {
		if resp.StatusCode != http.StatusOK {
			break
		}
		bytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("cannot read body: %s", err)
			break
		}
		pr := &persistableResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       string(bytes),
		}
		prJSON, err := json.MarshalIndent(pr, "", "  ")
		if err != nil {
			err = fmt.Errorf("cannot marshal JSON: %s", err)
			break
		}
		err = ioutil.WriteFile("../recorded/response.json", prJSON, 0644)
		if err != nil {
			err = fmt.Errorf("cannot write JSON: %s", err)
			break
		}
		hash := sha1.Sum(bytes)
		log.Printf("Response: %s", hex.EncodeToString(hash[:]))
	}
	return bytes, err
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
