package config

import "os"

type LoggerConfig struct {
	Enable             bool   `mapstructure:"enable" json:"enable"`
	EnableConsolePrint bool   `mapstructure:"enable_console_print" json:"enable_console_print"`
	EnableLogFiles     bool   `mapstructure:"enable_log_files" json:"enable_log_files"`
	FileLogFormat      string `mapstructure:"file_log_format" json:"file_log_format"`
	LogFilesPath       string `mapstructure:"log_files_path" json:"log_files_path"`
	LogFileName        string `mapstructure:"log_file_name" json:"log_file_name"`
	EnableTimeRolling  bool   `mapstructure:"enable_time_rolling" json:"enable_time_rolling"`
	EnableSizeRolling  bool   `mapstructure:"enable_size_rolling" json:"enable_size_rolling"`
	MaxFileSizeKB      int64  `mapstructure:"max_file_size_kb" json:"max_file_size_kb"`
	MaxFilesCount      int    `mapstructure:"max_files_count" json:"max_files_count"`
	Level              string `mapstructure:"level" json:"level"`
	Label              string `mapstructure:"label" json:"label"`
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Enable:             true,
		EnableConsolePrint: true,
		EnableLogFiles:     false,
		FileLogFormat:      "json",
		LogFilesPath:       ".",
		LogFileName:        "rzerolog.log",
		EnableTimeRolling:  false,
		EnableSizeRolling:  false,
		MaxFileSizeKB:      100 << 10,
		MaxFilesCount:      0,
		Level:              "DEBUG",
		Label:              "",
	}
}

func EnsureConfigPath(configPath string) error {
	return os.MkdirAll(configPath, os.ModePerm)
}
