package simple_mitm

const (
	ProxyPort           = "9999"
	APIServerPort       = "9998"
	APIServerAndPort    = "localhost:" + APIServerPort
	ModeKey             = "mode"
	CertVar             = "cert"
	KeyVar              = "key"
	CertDir             = "cert"
	GitHubAPIAndPort    = "api.github.com:443"
	CertExt             = "pem"
	CertFile            = CertDir + "/" + CertVar + "." + CertExt
	KeyFile             = CertDir + "/" + KeyVar + "." + CertExt
	HTTPSScheme         = "https"
	HTTPAuthHeader      = "Authorization"
	HTTPAuthBearerToken = "Bearer"
	//CACertFilepath      = CertDir + "/rootCA.crt"
	//GitHubAPIURL        = "https://" + GitHubAPIAndPort
)
