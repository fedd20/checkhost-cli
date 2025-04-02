package cmd

import (
	"checkhost-cli/logger"
	"checkhost-cli/utils"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Set configuration values",
	Long:  `You can set configuration values using the config command.`,
	Args:  cobra.MatchAll(cobra.MaximumNArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.GetLoggerInstance(true)

		config, err := utils.LoadConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Error loading config")
			return
		}

		all_settings := config.AllSettings()
		valid_data := map[string]any{}
		for key, value := range all_settings {
			if !strings.HasPrefix(key, "country_data") {
				valid_data[key] = value
			}
		}

		if len(args) == 0 {
			log.Info().Msg("Valid config:")
			for key, value := range valid_data {
				log.Info().Str("key", key).Str("value", fmt.Sprintf("%v", value)).Msg("-")
			}
		} else if args[0] == "list" {
			log.Info().Msg("Available keys:")
			for key := range valid_data {
				log.Info().Msgf("- %s", key)
			}
		} else if _, exists := valid_data[args[0]]; exists {
			if len(args) == 1 {
				value := config.GetString(args[0])
				log.Info().Str("key", args[0]).Str("value", value).Msg("Value from config")
			} else if len(args) == 2 {
				if b, err := strconv.ParseBool(args[1]); err == nil {
					config.Set(args[0], b)
				} else {
					if reflect.TypeOf(valid_data[args[0]]).Kind() == reflect.Bool {
						log.Fatal().Str("key", args[0]).Str("value", args[1]).Msg("Invalid value type")
					}
					config.Set(args[0], args[1])
				}

				err := config.WriteConfig()
				if err != nil {
					log.Fatal().Str("key", args[0]).Str("value", args[1]).Err(err).Msg("Error setting value")
					return
				}
				log.Info().Str("key", args[0]).Str("value", args[1]).Msg("Value set successfully!")
			}
		} else {
			log.Fatal().Str("key", args[0]).Msg("Key not found in config file")
		}
	},
}
