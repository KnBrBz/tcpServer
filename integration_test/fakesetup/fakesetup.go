package fakesetup

type Clients struct {
	host       string
	serverHost string
}

func NewClients(host, serverHost string) *Clients {
	return &Clients{
		host:       host,
		serverHost: serverHost,
	}
}

func (c *Clients) ServerHost() string {
	return c.serverHost
}

func (c *Clients) Host() string {
	return c.host
}

type Servers struct {
	host string
}

func NewServers(host string) *Servers {
	return &Servers{
		host: host,
	}
}

func (s *Servers) Host() string {
	return s.host
}
