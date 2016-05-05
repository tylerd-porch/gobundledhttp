package bundledhttp

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var pool *x509.CertPool

func init() {
	// Always build the pool
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts) // from certificates.go file
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
			DisableCompression: true,
		},
	}
}

// CtxBundled returns an oauth2 context with bundled http client
func CtxBundled() context.Context {
	return context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		bundledhttp.NewClient(),
	)
}

// Client without certificate checking.  Useful for self-signed certs.
func InsecureClient() *http.Client {
	// Insecure client without cert-trust checking
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
		},
	}

}

// Return just the certificate pool
func GetPool() *x509.CertPool {
	return pool
}
