package cmd

import (
	"fmt"
	"time"

	"github.com/SencilloDev/ngen/openapi"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openapiCmd = &cobra.Command{

	Use:          "openapi",
	Short:        "converts a nats micro info payload into OpenAPI spec",
	PreRunE:      validateConvertFlags,
	RunE:         convert,
	SilenceUsage: true,
}

var requiredConvertFlags = []string{
	"method_offset",
}

func init() {
	// attach convert subcommand to service subcommand
	generateCmd.AddCommand(openapiCmd)
	openapiFlags(openapiCmd)
	bindOpenapiFlags(openapiCmd)
}

func validateConvertFlags(cmd *cobra.Command, args []string) error {
	return validateEnvs(requiredConvertFlags...)
}

func convert(cmd *cobra.Command, args []string) error {
	nc, err := newNatsConnection("natsoapi")
	if err != nil {
		return err
	}
	defer nc.Close()

	name := viper.GetString("name")
	id := viper.GetString("id")
	methodOffset := viper.GetInt("method_offset")
	description := viper.GetString("description")
	title := viper.GetString("title")

	msg := nats.Msg{
		Subject: fmt.Sprintf("$SRV.INFO.%s.%s", name, id),
	}

	resp, err := nc.RequestMsg(&msg, 1*time.Second)
	if err != nil {
		return err
	}
	o := openapi.New(openapi.Opts{
		Version:      "3.0.0",
		Title:        title,
		Description:  description,
		MethodOffset: methodOffset,
	})

	spec, err := o.Convert(resp.Data)
	if err != nil {
		return err
	}

	fmt.Println(string(spec))

	return nil
}
