package config

// ExternalConfig will hold the configurations for external tools, such as Explorer or Elasticsearch
type ExternalConfig struct {
	ElasticSearchConnector ElasticSearchConfig
	EventNotifierConnector EventNotifierConfig
	WebSocketConnector     WebSocketDriverConfig
}

// ElasticSearchConfig will hold the configuration for the elastic search
type ElasticSearchConfig struct {
	Enabled                   bool
	IndexerCacheSize          int
	BulkRequestMaxSizeInBytes int
	URL                       string
	UseKibana                 bool
	Username                  string
	Password                  string
	EnabledIndexes            []string
}

// EventNotifierConfig will hold the configuration for the events notifier driver
type EventNotifierConfig struct {
	Enabled           bool
	UseAuthorization  bool
	ProxyUrl          string
	Username          string
	Password          string
	RequestTimeoutSec int
	MarshallerType    string
}

// CovalentConfig will hold the configurations for covalent indexer
type CovalentConfig struct {
	Enabled              bool
	URL                  string
	RouteSendData        string
	RouteAcknowledgeData string
}

// WebSocketDriverConfig will hold the configuration for WebSocket driver
type WebSocketDriverConfig struct {
	Enabled            bool
	IsServer           bool
	WithAcknowledge    bool
	BlockingAckOnError bool
	URL                string
	MarshallerType     string
	RetryDurationInSec int
}
