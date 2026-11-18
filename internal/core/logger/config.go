package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// этот файл помогает задать конфигурацию логгеру: какого уровня сообщения записывать и в какую папку сохранять логи

type LoggerConfig struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfig() (LoggerConfig, error) { // создаем и заполняем конфигурацию логгера
	var config LoggerConfig

	if err := envconfig.Process("LOGGER", &config); err != nil { // LOGGER - префикс
		return LoggerConfig{}, fmt.Errorf("process envconfig: %w", err) // у нас еще нет логгера, а ошибку надо как-то обработать
	}

	return config, nil
}

func NewLoggerConfigMust() LoggerConfig { // должен корректно появиться, иначе паника
	config, err := NewConfig()

	if err != nil { // так как у нас не получается создать конфигурацию для логгера и сам логгер соответственно, мы должны запаниковать
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}

	return config // возвращаем структуру, а не указатель так, как она используется один раз при создании логгера и больше не меняется
}
