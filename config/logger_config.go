// config/logger_config.go
package config

import "user-system/utils"

func SetupLogger() {
	utils.InitLogConfigs(map[string]utils.LogConfig{
		"app":   {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: false, AsJSON: false, Level: utils.INFO},
		"sql":   {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: true, AsJSON: false, Level: utils.INFO},
		"route": {ToConsole: true, ToFile: true, FilePath: "logs", IsDaily: true, AsJSON: false, Level: utils.INFO},
	})
}
