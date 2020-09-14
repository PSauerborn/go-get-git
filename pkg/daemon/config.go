package daemon

import (
    "os"
    "fmt"
    "strconv"
    log "github.com/sirupsen/logrus"
)

var (
    LogLevels = map[string]log.Level{ "DEBUG": log.DebugLevel, "INFO": log.InfoLevel, "WARN": log.WarnLevel }
    RabbitQueueUrl string
    QueueName string
    EventExchangeName string
    ExchangeType string
)

// Function used to configure service settings
func ConfigureService() {
    // set log level by overriding environment variables
    LogLevelString := OverrideStringVariable("LOG_LEVEL", "DEBUG")
    if LogLevel, ok := LogLevels[LogLevelString]; ok {
        log.SetLevel(LogLevel)
    } else {
        log.Fatal(fmt.Sprintf("received invalid log level %s", LogLevelString))
    }

    RabbitQueueUrl = OverrideStringVariable("GO_GET_GIT_RABBIT_QUEUE_URL", "amqp://guest:guest@192.168.99.100:5672/")
    QueueName = OverrideStringVariable("GO_GET_GIT_QUEUE_NAME", "testing-queue")
    EventExchangeName = OverrideStringVariable("GO_GET_GIT_EVENT_EXCHANGE_NAME", "events")
    ExchangeType = OverrideStringVariable("GO_GET_GIT_EVENT_EXCHANGE_TYPE", "fanout")
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideStringVariable(key string, DefaultValue string) string {
    value := os.Getenv(key)
    if len(value) > 0 {
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, value))
        return value
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideIntegerVariable(key string, DefaultValue int) int {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.Atoi(value)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to integer", result))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideFloatVariable(key string, DefaultValue float64) float64 {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.ParseFloat(value, 64)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to float", result))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}

// Function used to override configuration variables with some
// value by defaulting from environment variables
func OverrideBoolVariable(key string, DefaultValue bool) bool {
    value := os.Getenv(key)
    if len(value) > 0 {
        result, err := strconv.ParseBool(value)
        if err != nil {
            log.Fatal(fmt.Sprintf("cannot cast value '%v' to boolean", value))
        }
        log.Info(fmt.Sprintf("overriding variable %v with value %v", key, result))
        return result
    } else {
        return DefaultValue
    }
}