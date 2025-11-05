package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.PartsOrder = []string{"message"}

	logger := zerolog.New(output).With().Logger()

	var lines *int = pflag.IntP("tail", "n", 10, "Number of lines to show from the end of the logs for each container")
	pflag.Parse()
	if pflag.NArg() < 1 {
		pflag.Usage()
		os.Exit(1)
	}
	dc, err := NewDockerService()
	if err != nil {
		logger.Error().Err(err).Msg("Is docker service running?")
		panic(err)
	}

	containers, err := dc.GetContainers()
	if errors.Is(err, ErrNoContainers) {
		logger.Error().Msg("No containers found")
		return
	} else if err != nil {
		panic(err)
	}

	containerLogs := pflag.Args()

	containersToWatch := make([]string, 0)
	for key, name := range containers {
		for _, cLog := range containerLogs {
			if name == cLog {
				containersToWatch = append(containersToWatch, key)
				continue
			}
			if key[:len(cLog)] == cLog {
				containersToWatch = append(containersToWatch, key)
				continue
			}
		}
	}

	if len(containersToWatch) == 0 {
		logger.Error().Msg("No matching containers found")
		return
	}

	wg := sync.WaitGroup{}
	for i, id := range containersToWatch {
		wg.Add(1)
		go func() {

			err := dc.WatchLogs(id, containers[id], fmt.Sprintf("%d", *lines), Colors[i%len(Colors)])
			if err != nil {
				return
			}
		}()
	}

	awaitTermination(dc)
	wg.Wait()
}

func awaitTermination(dc *DockerService) {
	receiver := make(chan os.Signal)
	signal.Notify(receiver, os.Interrupt, os.Kill)

	<-receiver
	err := dc.Close()
	if err != nil {
		return
	}
	os.Exit(0)
}
