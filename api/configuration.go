package api

type ApiServerConfiguration struct {
	addr               string
	channelToOffServer <-chan bool
	newLogChannel      chan<- string
}

func NewApiServerConfiguration(addr string, channelToOffServer <-chan bool, newLogChannel chan<- string) *ApiServerConfiguration {
	return &ApiServerConfiguration{addr: addr, channelToOffServer: channelToOffServer, newLogChannel: newLogChannel}
}
