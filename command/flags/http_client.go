package flags

import (
	"flag"

	"github.com/hashicorp/consul/api"
)

type HTTPClient struct {
	Address       StringValue
	Token         StringValue
	CAFile        StringValue
	CAPath        StringValue
	CertFile      StringValue
	KeyFile       StringValue
	TLSServerName StringValue
}

func (f *HTTPClient) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(&f.Address, "http-addr",
		"The `address` and port of the Consul HTTP agent. The value can be an IP "+
			"address or DNS address, but it must also include the port. This can "+
			"also be specified via the CONSUL_HTTP_ADDR environment variable. The "+
			"default value is http://127.0.0.1:8500. The scheme can also be set to "+
			"HTTPS by setting the environment variable CONSUL_HTTP_SSL=true.")
	fs.Var(&f.Token, "token",
		"ACL token to use in the request. This can also be specified via the "+
			"CONSUL_HTTP_TOKEN environment variable. If unspecified, the query will "+
			"default to the token of the Consul agent at the HTTP address.")
	fs.Var(&f.CAFile, "ca-file",
		"Path to a CA file to use for TLS when communicating with Consul. This "+
			"can also be specified via the CONSUL_CACERT environment variable.")
	fs.Var(&f.CAPath, "ca-path",
		"Path to a directory of CA certificates to use for TLS when communicating "+
			"with Consul. This can also be specified via the CONSUL_CAPATH environment variable.")
	fs.Var(&f.CertFile, "client-cert",
		"Path to a client cert file to use for TLS when 'verify_incoming' is enabled. This "+
			"can also be specified via the CONSUL_CLIENT_CERT environment variable.")
	fs.Var(&f.KeyFile, "client-key",
		"Path to a client key file to use for TLS when 'verify_incoming' is enabled. This "+
			"can also be specified via the CONSUL_CLIENT_KEY environment variable.")
	fs.Var(&f.TLSServerName, "tls-server-name",
		"The server name to use as the SNI host when connecting via TLS. This "+
			"can also be specified via the CONSUL_TLS_SERVER_NAME environment variable.")
	return fs
}

// todo(fs): this should either go somewhere else or the pkg should not be named 'flags'
func NewAPIClient(f *HTTPClient) (*api.Client, error) {
	c := api.DefaultConfig()

	f.Address.Merge(&c.Address)
	f.Token.Merge(&c.Token)
	f.CAFile.Merge(&c.TLSConfig.CAFile)
	f.CAPath.Merge(&c.TLSConfig.CAPath)
	f.CertFile.Merge(&c.TLSConfig.CertFile)
	f.KeyFile.Merge(&c.TLSConfig.KeyFile)
	f.TLSServerName.Merge(&c.TLSConfig.Address)

	return api.NewClient(c)
}
