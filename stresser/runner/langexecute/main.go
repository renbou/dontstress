package main

import (
	"os"

	"github.com/renbou/dontstress/stresser/runner/langexecute/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal().Str("status", "invalid").Array("args", zerolog.Arr().Interface(os.Args)).Msg("Received less args than needed")
	}

	var (
		user string // user to run program as
		lang = os.Args[1]
		path = os.Args[2]
	)
	if user = os.Getenv("USER"); user == "" {
		log.Fatal().Str("status", "invalid").Msg("Must specify lowpriv user to run as in the env variable USER")
	}

	if valid, err := util.ValidFile(path); err != nil {
		log.Fatal().Str("status", "invalid").Msg("Unable to check " + path + " for validity")
	} else if valid == false {
		log.Fatal().Str("status", "invalid").Msg("File " + path + " does not exist or isn't a file")
	}

}
