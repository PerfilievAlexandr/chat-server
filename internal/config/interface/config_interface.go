package confog_interface

type GrpcServerConfig interface {
	Address() string
}

type DatabaseConfig interface {
	ConnectString() string
}
