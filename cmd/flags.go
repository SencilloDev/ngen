package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateEnvs(vals ...string) error {
	var errs []error
	for _, v := range vals {
		if !viper.IsSet(v) {
			errs = append(errs, fmt.Errorf("%v must be set", v))
		}
	}

	return errors.Join(errs...)
}

//Flags are defined here. Because of the way Viper binds values, if the same flag name is called
// with viper.BindPFlag multiple times during init() the value will be overwritten. For example if
// two subcommands each have a flag called name but they each have their own default values,
// viper can overwrite any value passed in for one subcommand with the default value of the other subcommand.
// The answer here is to not use init() and instead use something like PersistentPreRun to bind the
// viper values. Using init for the cobra flags is ok, they are only in here to limit duplication of names.

// bindNatsFlags binds nats flag values to viper
func bindNatsFlags(cmd *cobra.Command) {
	viper.BindPFlag("nats_urls", cmd.Flags().Lookup("nats-urls"))
	viper.BindPFlag("nats_seed", cmd.Flags().Lookup("nats-seed"))
	viper.BindPFlag("nats_jwt", cmd.Flags().Lookup("nats-jwt"))
	viper.BindPFlag("nats_secret", cmd.Flags().Lookup("nats-secret"))
	viper.BindPFlag("credentials_file", cmd.Flags().Lookup("credentials-file"))
}

// natsFlags adds the nats flags to the passed in cobra command
func natsFlags(cmd *cobra.Command) {
	cmd.Flags().String("nats-jwt", "", "NATS JWT as a string")
	cmd.Flags().String("nats-seed", "", "NATS seed as a string")
	cmd.Flags().String("credentials-file", "", "Path to NATS user credentials file")
	cmd.Flags().String("nats-urls", "nats://localhost:4222", "NATS URLs")
}

func convertFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "Name of the nats micro service")
	cmd.Flags().StringP("id", "i", "", "ID of the NATS micro service")
	cmd.Flags().StringP("title", "t", "", "Title for OpenAPI spec")
	cmd.Flags().StringP("description", "d", "", "Description for OpenAPI spec")
	cmd.Flags().IntP("method-offset", "m", 0, "Offset of the subject to get the HTTP method")
}

func bindConvertFlags(cmd *cobra.Command) {
	viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("id", cmd.Flags().Lookup("id"))
	viper.BindPFlag("title", cmd.Flags().Lookup("title"))
	viper.BindPFlag("description", cmd.Flags().Lookup("description"))
	viper.BindPFlag("method_offset", cmd.Flags().Lookup("method-offset"))
}
