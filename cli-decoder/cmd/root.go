package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"workshops/cli-decoder/internal"
)

//TODO: move to env
const StoragePath = "assets"

var (
	rootCmd = &cobra.Command{
		Use:   "cli-decoder",
		Short: "A decoder for xml and json to file",
		Long:  `A decoder for xml and json and save unique strings into files`,
		Run: func(cmd *cobra.Command, args []string) {
			scanInput()
		},
	}
	jsonFlag, xmlFlag bool
	decoder           string
	handler           Handler
)

type Handler struct {
	d internal.CliDecoder
	s internal.Storage
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Use json decoder")
	rootCmd.PersistentFlags().BoolVarP(&xmlFlag, "xml", "x", false, "Use xml decoder")
}

func initConfig() {
	var decoderType string
	/*
		On startup identify all known md5 hashsums
		Put hash data into the separate file and keep them in memory when program is running
	*/
	if jsonFlag {
		decoderType = internal.JSON
	} else if xmlFlag {
		decoderType = internal.XML
	} else {
		log.Fatal("--json or --xml flag is required")
	}
	handler = Handler{
		s: internal.NewDataStorage(StoragePath, decoderType),
		d: internal.NewCliDecoder(decoderType),
	}
}

func scanInput() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	// when you hit enter - the data is read
	for scanner.Scan() {
		yourData := scanner.Bytes()

		err := handler.d.Decode(yourData)
		if err != nil {
			log.Printf("Cant decode your data as %s", decoder)
			continue
		}
		err = handler.s.Save(yourData)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
