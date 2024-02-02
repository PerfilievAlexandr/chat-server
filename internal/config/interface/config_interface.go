package confog_interface

type GrpcAuthClientConfig interface {
	Address() string
}

type GrpcServerConfig interface {
	Address() string
}

type DatabaseConfig interface {
	ConnectString() string
}
