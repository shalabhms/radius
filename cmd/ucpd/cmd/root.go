// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/project-radius/radius/pkg/trace"
	"github.com/project-radius/radius/pkg/ucp/server"
	"github.com/project-radius/radius/pkg/ucp/ucplog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ucpd",
	Short: "UCP server",
	Long:  `Server process for the Univeral Control Plane (UCP).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := server.NewServerOptionsFromEnvironment()
		if err != nil {
			return err
		}

		logger, flush, err := ucplog.NewLogger(ucplog.LoggerName, &options.LoggingOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer flush()

		host, err := server.NewServer(options)
		if err != nil {
			return err
		}
		ctx := logr.NewContext(cmd.Context(), logger)
		ctx, cancel := context.WithCancel(ctx)

		options.TracerProviderOptions.ServiceName = server.ServiceName
		shutdown, err := trace.InitTracer(options.TracerProviderOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal("failed to shutdown TracerProvider: %w", err)
			}
		}()

		stopped, serviceErrors := host.RunAsync(ctx)

		exitCh := make(chan os.Signal, 2)
		signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

		select {
		// Shutdown triggered
		case <-exitCh:
			logger.Info("Shutting down....")
			cancel()

		// A service terminated with a failure. Shut down
		case <-serviceErrors:
			logger.Info("Error occurred - shutting down....")
			cancel()
		}

		// Finished shutting down. An error returned here is a failure to terminate
		// gracefully, so just crash if that happens.
		err = <-stopped
		if err == nil {
			os.Exit(0)
		} else {
			panic(err)
		}

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.ExecuteContext(context.Background()))
}
