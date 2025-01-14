package cmd

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jinx/src/jinkiesengine"
	"os"
)

var (
	containerConfigPath = ""
	hostConfigPath      = ""
	hostConfig          container.HostConfig
)

func hydrateFromConfig(configPath string) jinkiesengine.ContainerInfo {
	var config jinkiesengine.ContainerInfo

	if configPath == "" {
		config.ImageName = "jamandbees/jinkies"
		config.ContainerName = "jinkies"
		config.ContainerPort = "8080/tcp"
		config.HostIp = "0.0.0.0"
		config.HostPort = "8090/tcp"
		config.PullImages = true
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigType("env")
		viper.SetConfigName(configPath)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
		viper.Unmarshal(&config)
	}

	return config
}

func addConfig(configPath string) container.HostConfig {
	var config container.HostConfig

	viper.AddConfigPath("./")
	viper.SetConfigType("yml")
	viper.SetConfigName(configPath)

	if err := viper.ReadInConfig(); err != nil {
		config = container.HostConfig{
			AutoRemove: true,
		}
	}
	viper.Unmarshal(&config)

	return config
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Subcommands to allow you to start or stop an unconfigured jinkies",
	Long: `Why would you want an unconfigured instance of jinkies? Any time you want a jenkins instance
quickly for reasons unrelated to a specific job. Maybe you want to prototype some jcasc settings or something.

Maybe you want two instances of jinkies running at once? Use the -c flag to supply an environment variables file. To
write a blank version of this file, see the 'jinx containerconfig' subcommand.
`,
}

var startSubCmd = &cobra.Command{
	Use:   "start",
	Short: "start jinkies!",
	Long:  `Starts the unconfigured jinkies container`,
	Run: func(cmd *cobra.Command, args []string) {
		jinkiesengine.RunRunRun(hydrateFromConfig(containerConfigPath), addConfig(hostConfigPath))
	},
}

var stopSubCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops your jinkies container_info.",
	Long:  `No configuration is retained after a stop, so this gets you back to a clean slate.`,
	Run: func(cmd *cobra.Command, args []string) {
		jinkiesengine.StopGirl(hydrateFromConfig(containerConfigPath))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.AddCommand(startSubCmd)
	serveCmd.AddCommand(stopSubCmd)

	serveCmd.PersistentFlags().StringVarP(&containerConfigPath, "containerconfig", "c", "", "Path to config file describing your container")
	serveCmd.PersistentFlags().StringVarP(&hostConfigPath, "hostconfig", "o", "", "Path to config file describing your container host ")

	if hostConfigPath != "" {
		_, err := os.Open(hostConfigPath)
		if err != nil {
			fmt.Printf("Could not open host config file %v \n", err)
		}
	}

	if containerConfigPath != "" {
		_, error := os.Open(containerConfigPath)
		if error != nil {
			fmt.Printf("Could not open container config file %v \n", error)
		}
	}

}
