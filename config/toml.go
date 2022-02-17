package config

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/spf13/viper"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

#####################
### Logger params ###
#####################
# Enable logger
enable = {{ .Enable}}
# Whether print log text on console
enable_console_print = {{ .EnableConsolePrint}}
# Whether output log to files
enable_log_files = {{ .EnableLogFiles}}
# Output format to log files
# ["text","json"] supported
file_log_format = "{{ .FileLogFormat}}"
# Whether enable time rolling rules
# If true set, the final log file name will be parsed by time parser with 'log_file_name'
#   format:
#       yyyy - year
#       MM   - month
#       dd   - day
#       HH   - hour
#       mm   - minute
#       ss   - second
#   eg: "rzerolog-yyyyMMddHH.log" -> "rzerolog-2022021510.log"
#   also support golang format just like "rzerolog-200601021504.log" -> "rzerolog-202202151055.log"
enable_time_rolling = {{ .EnableTimeRolling}}
# Whether enable size rolling rules
enable_size_rolling = {{ .EnableSizeRolling}}
# Max size in Kb of each log file
# If size of file reach the value, the file will be added a suffix such as '.1'
max_file_size_kb = {{ .MaxFileSizeKB}}
# Max count of log files saved
# Files too old will be cleared
max_files_count = {{ .MaxFilesCount}}
# Path of log files
log_files_path = "{{ .LogFilesPath}}"
# Name of log files
log_file_name = "{{ .LogFileName}}"
# Logger level
level = "{{ .Level}}"
# Logger label
label = "{{ .Label}}"
`

func WriteConfigToTomlFile(configFilePath string, config *LoggerConfig) error {
	var buffer bytes.Buffer
	tmpl := template.New("loggerConfigFileTemplate")
	configTemplate, err := tmpl.Parse(defaultConfigTemplate)
	if err != nil {
		return err
	}

	if err = configTemplate.Execute(&buffer, config); err != nil {
		return err
	}

	return ioutil.WriteFile(configFilePath, buffer.Bytes(), 0666)
}

func GetLoggerConfigFromFile(configFilePath string, v *viper.Viper) (*LoggerConfig, error) {
	if v == nil {
		v = viper.New()
	}
	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := new(LoggerConfig)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetLoggerConfigFromPath(configPath string, v *viper.Viper) (*LoggerConfig, error) {
	if v == nil {
		v = viper.New()
	}
	v.AddConfigPath(configPath)
	v.SetConfigName("logger")
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := new(LoggerConfig)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
