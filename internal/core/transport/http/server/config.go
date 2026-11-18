package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"` // адрес самого сервера
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

func NewHTTPConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("HTTP", &config); err != nil { // HTTP - префикс
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewHTTPConfigMust() Config { // должен корректно появиться, иначе паника
	config, err := NewHTTPConfig()

	if err != nil { // так как у нас не получается создать конфигурацию для HTTP-сервера и сам HTTP-сервер соответственно, мы должны запаниковать
		err = fmt.Errorf("get HTTP config: %w", err)
		panic(err)
	}

	return config // возвращаем структуру, а не указатель так, как она используется один раз при создании HTTP-сервера и больше не меняется
}
