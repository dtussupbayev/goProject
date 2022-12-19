package message_broker

type CacheBroker interface {
	BrokerWithClient
	Remove(key interface{}) error //ubrat po key
	Purge() error                 //voobshe ochistit cache
}
