package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/SencilloDev/ngen/diagram"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Generate a d2 diagram of subject mapping",
	RunE:  generateDiagram,
}

func init() {
	generateCmd.AddCommand(diagramCmd)
	diagramCmd.Flags().StringP("out", "o", "diagram.svg", "Output file name")
	viper.BindPFlag("out", diagramCmd.Flags().Lookup("out"))
	diagramCmd.Flags().BoolP("print", "p", false, "Print diagram to stdout")
	viper.BindPFlag("print", diagramCmd.Flags().Lookup("print"))
	diagramCmd.Flags().BoolP("animate", "a", false, "Animate edges")
	viper.BindPFlag("animate", diagramCmd.Flags().Lookup("animate"))
}

func generateDiagram(cmd *cobra.Command, args []string) error {
	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	nc, err := newNatsConnection("natsdiagrams")
	if err != nil {
		return err
	}
	defer nc.Close()

	name := viper.GetString("name")
	id := viper.GetString("id")

	msg := nats.Msg{
		Subject: fmt.Sprintf("$SRV.INFO.%s.%s", name, id),
	}

	resp, err := nc.RequestMsg(&msg, 1*time.Second)
	if err != nil {
		return err
	}

	var m micro.Info
	if err := json.Unmarshal(resp.Data, &m); err != nil {
		return err
	}

	opts := diagram.GraphOpts{
		// if printing set generate to false
		GenerateSVG: !viper.GetBool("print"),
	}

	if viper.GetBool("animate") {
		opts.EdgeOpts = append(opts.EdgeOpts, diagram.WithAnimation)
	}

	text, data, err := diagram.New(cmd.Context(), logger, m, opts)
	if err != nil {
		return err
	}

	if viper.GetBool("print") {
		fmt.Println(text)
		return nil
	}

	out := viper.GetString("out")
	return os.WriteFile(filepath.Join(out), data, 0600)
}
