package etcd

import "mxshop-srvs/inventory_srv/utils/registry/consul"

type client struct{}

func NewClient() consul.RegisterClient {
	return &client{}
}

func (client) Register(srvID, name string, tags []string, port int, addr string) error {
	panic("implement me")
}

func (client) Deregister(srvID string) error {
	panic("implement me")
}
