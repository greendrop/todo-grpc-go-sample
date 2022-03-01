package model

type AppConfig struct {
	AppEnv   string
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Url string
}

type ServerConfig struct {
	GrpcPort        string `toml:"grpc_port"`
	GrpcGatewayPort string `toml:"grpc_gateway_port"`
}
