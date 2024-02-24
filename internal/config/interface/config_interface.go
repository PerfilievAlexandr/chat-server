package configInterface

type GrpcAuthClientConfig interface {
	Address() string
}

type GrpcServerConfig interface {
	Address() string
}

type DatabaseConfig interface {
	ConnectString() string
}

type PrometheusServerConfig interface {
	Address() string
}
