/*
Config собирает конфигурацию приложения.
Наибольший приоритет отдаётся значениям параметров, установленным в флагах при запуске приложения.
Если значение параметра не установлено флагом, его значение ищется в перемнных окружения.
Если значение параметра не найдено, используются значения по умолчанию, установленные константами.
*/
package config

import (
	"flag"
	"os"

	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

type TypeRepo int

const (
	Mem TypeRepo = iota
	PG
)

const (
	FlagHTTPPort = "port"      // Имя флага для порта HTTP-сервера
	EnvHTTPPort  = "HTTP_PORT" // Название переменной среды, в которой ищем порт, на котором будет работать HTTP-сервер
	DefHTTPPort  = "8088"      // Значение порта HTTP-сервера по умолчанию

	FlagRPCHost = "hostRPC"   // Имя флага для хоста gRPC
	EnvRPCHost  = "RPC_HOST"  // Название переменной среды, в которой ищем хост, на котором ищем gRPC
	DefRPCHost  = "localhost" // Значение хоста gRPC по умолчанию

	FlagRPCPort = "portRPC"  // Имя флага для порта gRPC
	EnvRPCPort  = "RPC_PORT" // Название переменной среды, в которой ищем порт, на котором ищем gRPC
	DefRPCPort  = "50051"    // Значение порта gRPC по умолчанию
)

type (
	// Config - конфгурация приложения
	Config struct {
		ConfigHTTP
		repo_client.ConfigRPC
	}

	// ConfigHTTP - конфигурация HTTP сервера
	ConfigHTTP struct {
		Port string
	}
)

// config.New() возвращает ссылку на конфиг и ошибку-результат инициализации конфига
func New() (*Config, error) {
	cfg := &Config{}
	var err error
	var ok bool

	zap.S().Debug("Конфигурация приложения\n")
	// порт http сервера
	var port string
	if port, ok = os.LookupEnv(EnvHTTPPort); !ok {
		port = DefHTTPPort
	}

	/*
		в StringVar &port - куда запишется итоговое значение
		если найдёт флаг, то его, иначе - значение по умолчанию = port,
		то есть сохранит старое значение
	*/
	flag.StringVar(&port, FlagHTTPPort, port, "usage: b2b2b -port=$NNNN")

	var hostRPC, portRPC string
	if hostRPC, ok = os.LookupEnv(EnvRPCHost); !ok {
		hostRPC = DefRPCHost
	}

	if portRPC, ok = os.LookupEnv(EnvRPCPort); !ok {
		portRPC = DefRPCPort
	}

	flag.StringVar(&hostRPC, FlagRPCHost, hostRPC, "usage: b2b2b -hostRPC=$NNNN")
	flag.StringVar(&portRPC, FlagRPCPort, portRPC, "usage: b2b2b -portRPC=$NNNN")

	flag.Parse()

	cfg.ConfigHTTP.Port = port
	cfg.ConfigRPC.HostRPC = hostRPC
	cfg.ConfigRPC.PortRPC = portRPC

	zap.S().Debug("HTTP Port: ", cfg.ConfigHTTP.Port, "\n")

	return cfg, err
}
