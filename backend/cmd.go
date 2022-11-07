package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// ... extension of the config file that we are going to append to the file if we are using
	// ${PWD} or ${HOME} variables instead of config-file parameter.
	configFileType = "yml"
	// ... name of the config file that we are going to prepend to the file if we are using
	// ${PWD} or ${HOME} variable instead of config-file parameter.
	configFileName = "main"
)

var (
	// ... can't check the file because it doesn't exists, or we don't have permissions.
	errCheckConfigFile = errors.New("checking config file")

	// ... the name of the config file
	configFile string
	// ... the global configuration
	cfg Config
)

// ... newRootCmd allows us to create the main command to configure the application
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "backend",
		Short: "backend of the application",
		Long:  "backend of the application",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile != "" {
				readConfig()
			}
		},
	}

	root.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "config file to parse information. Command line args are preferred over config file ones")

	return root
}

// ... we try to read the config file and update the global config structure
func readConfig() error {
	if !checkConfigFile() {
		return errCheckConfigFile
	}

	if viper.GetViper().ConfigFileUsed() == "" {
		viper.SetConfigType(configFileType)
		viper.SetConfigName(configFileName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	ConfigureLogger(cfg.Logger)
	configureDB(cfg.Database)

	return nil
}

// ... we try to read first the config-file parameter if it is not empty. If we can't find the file or it doesn't exist, we use some default paths
// like ${PWD} or ${HOME}, taking preference actual directory, aka ${PWD}
func checkConfigFile() bool {
	var exists bool = false

	if configFile != "" && fileExists(configFile) {
		viper.SetConfigFile(configFile)
		return true
	}

	if pwd, err := os.Getwd(); err == nil && fileExists(buildCfg(pwd)) {
		exists = true
		viper.AddConfigPath(pwd)
	}

	if home, err := os.UserHomeDir(); err == nil && fileExists(buildCfg(home)) {
		exists = true
		viper.AddConfigPath(home)
	}

	return exists
}

func fileExists(file string) bool {
	_, err := os.OpenFile(file, os.O_RDONLY, 0444)
	return err == nil
}

// ... name of the configuration by default that we have to use
func buildCfg(dir string) string {
	return filepath.Join(dir, configFileName+"."+configFileType)
}

// ... main entrypoint of the cmd file, so we can configure the app and observe if there is something wrong
func execute() error {
	root := newRootCmd()

	readConfig() //nolint:errcheck

	err := root.Execute()
	if err != nil {
		return err
	}

	return nil
}
