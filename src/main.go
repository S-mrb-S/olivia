package main

/*
   #include "stdio.h"
   void wrapPrintf(const char *s) {
      printf("%s", s);
   }
*/
import "C"
import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"

	global "github.com/S-mrb-S/olivia/core"
)

var neuralNetworksMapContainer = map[string]global.Network{}
const (
	defaultPort    = "2006"
)

func main() {
	C.wrapPrintf(C.CString("Hello, World\n"))

	serverPortArg := flag.String("port", defaultPort, "The port for the API and WebSocket.")
	localeRetrainArg := flag.String("re-train", "", "The locale(s) to re-train.")
	flag.Parse()

	// If the localeRetrainArg isn't empty then retrain the given models
	if *localeRetrainArg != "" {
		executeModelRetraining(*localeRetrainArg)
	}

	// Print the Olivia ASCII text
	oliviaASCIIBanner := string(global.FetchFileContent("../res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCIIBanner))

	// Create the authentication token
	global.Authenticate()

	for _, individualLocale := range global.Locales {
		global.GenerateSerializedMessages(individualLocale.Tag)

		neuralNetworksMapContainer[individualLocale.Tag] = global.CreateNeuralNetwork(
			individualLocale.Tag,
			false,
		)
	}

	// Get port from environment variables if there is
	if os.Getenv("PORT") != "" {
		*serverPortArg = os.Getenv("PORT")
	}

	// Serves the server
	global.StartServer(neuralNetworksMapContainer, *serverPortArg)
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
