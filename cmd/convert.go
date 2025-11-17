package cmd

import (
	"fmt"
	"time"

	"github.com/SencilloDev/natsoapi/openapi"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertCmd = &cobra.Command{
	Use:          "convert",
	Short:        "convert converts a nats micro info payload into OpenAPI spec",
	PreRunE:      validateConvertFlags,
	RunE:         convert,
	SilenceUsage: true,
}

var requiredConvertFlags = []string{
	"name",
	"id",
	"method_offset",
}

func init() {
	// attach convert subcommand to service subcommand
	rootCmd.AddCommand(convertCmd)
	natsFlags(convertCmd)
	bindNatsFlags(convertCmd)
	convertFlags(convertCmd)
	bindConvertFlags(convertCmd)
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
