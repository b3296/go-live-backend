// config/logger_config.go
package config

import "user-system/utils"

func SetupLogger() {
	utils.InitLogConfigs(map[string]utils.LogConfig{
		"app":   {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: false},
		"sql":   {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: true},
		"route": {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: true},
	})
}
