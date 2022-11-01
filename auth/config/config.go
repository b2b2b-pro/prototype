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

	"go.uber.org/zap"
)

const (
	FlagHTTPPort = "port"      // Имя флага для порта HTTP-сервера
	EnvHTTPPort  = "HTTP_PORT" // Название переменной среды, в которой ищем порт, на котором будет работать HTTP-сервер
	DefHTTPPort  = "8888"      // Значение порта HTTP-сервера по умолчанию
)

type (
	// Config - конфгурация приложения
	Config struct {
		ConfigHTTP
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

	flag.Parse()

	cfg.ConfigHTTP.Port = port

	zap.S().Debug("HTTP Port: ", cfg.ConfigHTTP.Port, "\n")

	return cfg, err
}
