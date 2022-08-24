package util

import "github.com/spf13/viper"

// All configuration for application
type Config struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	SendGridKey  string `mapstructure:"SENDGRID_API_KEY"`
	SlackIssueID string `mapstructure:"SLACK_ISSUE_ID"`
	SlackErrorID string `mapstructure:"SLACK_ERROR_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
