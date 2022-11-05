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

	"github.com/b2b2b-pro/lib/repo_srv"
	"go.uber.org/zap"
)

/*
github.com/b2b2b-pro/prototype/repository
*/
const (
	FlagPGHost = "hostPG"    // Имя флага для хоста СУБД
	EnvPGHost  = "PG_HOST"   // Название переменной среды, в которой ищем хост, на котором ищем СУБД
	DefPGHost  = "localhost" // Значение хоста СУБД по умолчанию

	FlagPGPort = "portPG"  // Имя флага для порта СУБД
	EnvPGPort  = "PG_PORT" // Название переменной среды, в которой ищем порт, на котором ищем СУБД
	DefPGPort  = "5432"    // Значение порта СУБД по умолчанию

	FlagDBName = "DBName"  // Имя флага имени СУБД
	EnvDBName  = "DB_NAME" // Название переменной среды, в которой ищем имя СУБД
	DefDBName  = "b2b2b"   // Значение имени СУБД по умолчанию

	FlagDBUser = "DBUser"  // Имя флага для пользователя СУБД
	EnvDBUser  = "DB_USER" // Название переменной среды, в которой пользователя СУБД
	DefDBUser  = "user"    // Значение пользователяа СУБД по умолчанию

	FlagDBPass = "DBUPass" // Имя флага для пользователя СУБД
	EnvDBPass  = "DB_Pass" // Название переменной среды, в которой пользователя СУБД
	DefDBPass  = "pass"    // Значение пользователяа СУБД по умолчанию

	FlagRPCHost = "hostRPC"  // Имя флага для хоста gRPC
	EnvRPCHost  = "RPC_HOST" // Название переменной среды, в которой ищем хост, на котором поднимаем gRPC
	DefRPCHost  = ""         // Значение хоста gRPC по умолчанию

	FlagRPCPort = "portRPC"  // Имя флага для порта gRPC
	EnvRPCPort  = "RPC_PORT" // Название переменной среды, в которой ищем порт, на котором поднимаем gRPC
	DefRPCPort  = "50051"    // Значение порта gRPC по умолчанию
)

type (
	// Config - конфгурация приложения
	Config struct {
		ConfigPG
		repo_srv.ConfigRPC
	}

	// ConfigPG - конфигурация postgresql
	ConfigPG struct {
		HostPG string
		PortPG string
		DBName string
		DBUser string
		DBPass string
	}
)

// TODO поискать более красивый способ сбора конфига
// config.New() возвращает ссылку на конфиг и ошибку-результат инициализации конфига
func New() (*Config, error) {
	cfg := &Config{}
	var err error
	var ok bool

	zap.S().Debug("Конфигурация приложения\n")

	var hostPG, portPG, dbName, dbUser, dbPass string
	if hostPG, ok = os.LookupEnv(EnvPGHost); !ok {
		hostPG = DefPGHost
	}

	if portPG, ok = os.LookupEnv(EnvPGPort); !ok {
		portPG = DefPGPort
	}

	if dbName, ok = os.LookupEnv(EnvDBName); !ok {
		dbName = DefDBName
	}

	if dbUser, ok = os.LookupEnv(EnvDBUser); !ok {
		dbUser = DefDBUser
	}

	if dbPass, ok = os.LookupEnv(EnvDBPass); !ok {
		dbPass = DefDBPass
	}

	/*
		в StringVar &portPG - куда запишется итоговое значение
		если найдёт флаг, то его, иначе - значение по умолчанию = portPG,
		то есть сохранит старое значение
	*/
	flag.StringVar(&hostPG, FlagPGHost, hostPG, "usage: b2b2b -hostPG=$NNNN")
	flag.StringVar(&portPG, FlagPGPort, portPG, "usage: b2b2b -portPG=$NNNN")
	flag.StringVar(&dbName, FlagDBName, dbName, "usage: b2b2b -DBName=$NNNN")
	flag.StringVar(&dbUser, FlagDBUser, dbUser, "usage: b2b2b -DBUser=$NNNN")
	flag.StringVar(&dbPass, FlagDBPass, dbPass, "usage: b2b2b -DBPass=$NNNN")

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

	cfg.ConfigPG.HostPG = hostPG
	cfg.ConfigPG.PortPG = portPG
	cfg.ConfigPG.DBName = dbName
	cfg.ConfigPG.DBUser = dbUser
	cfg.ConfigPG.DBPass = dbPass

	cfg.ConfigRPC.HostRPC = hostRPC
	cfg.ConfigRPC.PortRPC = portRPC

	return cfg, err
}
