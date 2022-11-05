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

/*
type TypeRepo int

const (
	Mem TypeRepo = iota
	PG
)
*/

const (
	FlagRPCHost = "hostRPC"   // Имя флага для хоста gRPC
	EnvRPCHost  = "RPC_HOST"  // Название переменной среды, в которой ищем хост, на котором ищем gRPC
	DefRPCHost  = "localhost" // Значение хоста gRPC по умолчанию

	FlagRPCPort = "portRPC"  // Имя флага для порта gRPC
	EnvRPCPort  = "RPC_PORT" // Название переменной среды, в которой ищем порт, на котором ищем gRPC
	DefRPCPort  = "50051"    // Значение порта gRPC по умолчанию
)

// Config - конфгурация приложения
type Config struct {
	ConfigRPC repo_client.ConfigRPC
	Tkn       string
}

// config.New() возвращает ссылку на конфиг и ошибку-результат инициализации конфига
func New() (*Config, error) {

	// TODO пока прибил гвоздями долгий login admin'а
	cfg := &Config{
		Tkn: "zlt/beUeYzGOevstEzcuZIvimiNDwuZpTRPBw/LpPan6WySdHeAol/uP4pJY0M/QinILP5AzrPhKGVQT1diOBij4pMXm9iVXmc1HJQIeLcjBNFQiJZokYH1nFHk4n2/4br3nmUqE6zK9kJTO8kCLOb1TJMv6LRQR2sk2QCZVth9ZUE9d2pB2ckMc/+xBNY8XwnkGd/hvNiNJjny79mwcF8cjECbaQ8P53mDLtS99qp8akk6mDsRg21k0hmmC+2Z4XOOCFjqX7B0V2/5byfaXoLS3AqnjkJ87nIaemQk5wrh0aDqMtffU4+1yqVuQhpLvAxa8UWdF6Npv0oZsVxIb4SeL/AbZ/MXZWJRtR4OOXCrfp/fjh2VGyOOo2IE9esw8j23RhVLdBwRRjKpWoA==",
	}
	var err error
	var ok bool

	zap.S().Debug("Конфигурация приложения\n")

	/*
		в StringVar &port - куда запишется итоговое значение
		если найдёт флаг, то его, иначе - значение по умолчанию = port,
		то есть сохранит старое значение
	*/

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

	cfg.ConfigRPC.HostRPC = hostRPC
	cfg.ConfigRPC.PortRPC = portRPC

	zap.S().Debugf("Config: %v\n", cfg)

	return cfg, err
}
