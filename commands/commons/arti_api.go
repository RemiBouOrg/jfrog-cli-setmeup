package commons

type ProxySettings struct {
	ServerName               string
	UseHttp                  bool
	UseHttps                 bool
	HttpPort                 int
	HttpsPort                int
	DockerReverseProxyMethod string
	ReverseProxyRepositories ReverseProxyRepositories
}

type ReverseProxyRepositories struct {
	ReverseProxyRepoConfigs []ReverseProxyRepoConfigs
}

type ReverseProxyRepoConfigs struct {
	RepoRef    string
	Port       int
	ServerName string
}
