package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"

	olivia "github.com/S-mrb-S/olivia/core"
)

var neuralNetworksMapContainer = map[string]olivia.Network{}
const (
	defaultPort    = "2006"
)

func main() {
	serverPortArg := flag.String("port", defaultPort, "The port for the API and WebSocket.")
	localeRetrainArg := flag.String("re-train", "", "The locale(s) to re-train.")
	flag.Parse()

	// If the localeRetrainArg isn't empty then retrain the given models
	if *localeRetrainArg != "" {
		executeModelRetraining(*localeRetrainArg)
	}

	// Print the Olivia ASCII text
	oliviaASCIIBanner := string(olivia.FetchFileContent("../res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCIIBanner))

	// Create the authentication token
	olivia.Authenticate()

	for _, individualLocale := range olivia.Locales {
		olivia.GenerateSerializedMessages(individualLocale.Tag)

		neuralNetworksMapContainer[individualLocale.Tag] = olivia.CreateNeuralNetwork(
			individualLocale.Tag,
			false,
		)
	}

	// Serves the server
	olivia.StartServer(neuralNetworksMapContainer, *serverPortArg)
}

// executeModelRetraining retrains the given locales
func executeModelRetraining(localeRetrainList string) {
	// Iterate locales by separating them by comma
	for _, individualLocale := range strings.Split(localeRetrainList, ",") {
		trainingFilePath := fmt.Sprintf("../res/locales/%s/training.json", individualLocale)
		deleteError := os.Remove(trainingFilePath)

		if deleteError != nil {
			fmt.Printf("Cannot re-train %s model.", individualLocale)
			return
		}
	}
}
