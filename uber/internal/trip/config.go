package trip

import "time"

type AppConfig struct {
	RESTBindAddress string `envconfig:"REST_BIND_ADDRESS" default:":8080"`
	GRPCBindAddress string `envconfig:"GRPC_BIND_ADDRESS" default:":9090"`

	Mongo MongoConfig `envconfig:"MONGO"`
}

type MongoConfig struct {
	URI            string        `envconfig:"URI" default:"mongodb://localhost:27017"`
	Database       string        `envconfig:"DATABASE" default:"location"`
	MaxPoolSize    uint64        `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize    uint64        `envconfig:"MIN_POOL_SIZE" default:"10"`
	ConnectTimeout time.Duration `envconfig:"CONNECT_TIMEOUT" default:"10s"`
}

type RedisConfig struct {
	Address  string `envconfig:"ADDRESS" default:"localhost:6379"`
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	DB       int    `envconfig:"DB" default:"0"`
}
