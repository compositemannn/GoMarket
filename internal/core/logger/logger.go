package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger // встраивание структуры в собственную - все методы доступны для Logger

	file *os.File // файл для логов
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context") // если логгер не получилось достать из контекста, значит его там нет
	}
	return log
}

func NewLogger(loggerConfig LoggerConfig) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel() // уровень можно будет безопасно менять при работе нескольких горутин
	if err := zapLvl.UnmarshalText([]byte(loggerConfig.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}
	if err := os.MkdirAll(loggerConfig.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000") // имя файла по дате
	logFilePath := filepath.Join(                                      // соединяет путь к папке  именем файла
		loggerConfig.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644) // открывает файл и создает при отсутствии
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig) // знает как форматировать сообщения
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl), // вывод логов в терминал
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),   // вывод логов в файл
	)

	zapLogger := zap.New(core, zap.AddCaller()) // AddCaller - доп. инфа о месте вызова: исходный файл и строка

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger: ", err)
	}
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}
