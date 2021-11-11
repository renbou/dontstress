package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/renbou/dontstress/stresser/runner/langexecute/runner"
	"github.com/renbou/dontstress/stresser/runner/langexecute/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LAUNCH_SINGLE = "single"
	LAUNCH_TEST   = "test"
)

func FatalInternal() *zerolog.Event {
	return log.Fatal().Str("type", "internal")
}

func FatalUser() *zerolog.Event {
	return log.Fatal().Str("type", "user")
}

func HandlePrepareError(err error) {
	re := &runner.RunError{}
	if err != nil && errors.As(err, &re) {
		if re.Step() == runner.STEP_INIT {
			FatalInternal().Msg("Runner init failed: " + re.Error())
		} else if re.Step() == runner.STEP_COMPILE {
			FatalUser().Msg("Compilation Error: " + re.Unwrap().Error())
		} else {
			FatalInternal().Msg("Unknown RunError from Prepare: " + re.Error())
		}
	} else if err != nil {
		FatalInternal().Msg("Unknown error from Prepare: " + err.Error())
	}
}

func HandleRunError(err error) {
	re := &runner.RunError{}
	if err != nil && errors.As(err, &re) {
		if re.Step() == runner.STEP_RUN {
			FatalUser().Msg("Runtime Error: " + re.Error())
		} else {
			FatalInternal().Msg("Unknown RunError from Run: " + re.Error())
		}
	} else if err != nil {
		FatalInternal().Msg("Unknown error from Run: " + err.Error())
	}
}

func Run(run runner.Runner, input string) string {
	b := new(strings.Builder)
	HandleRunError(run.Run(strings.NewReader(input), b))
	return b.String()
}

func main() {
	if len(os.Args) < 4 {
		FatalInternal().Array("args", zerolog.Arr().Interface(os.Args)).Msg("Received less args than needed")
	}

	var (
		//user       string // user to run program as
		lang       = os.Args[1]
		path       = os.Args[2]
		launchType = os.Args[3]
	)
	// if user = os.Getenv("USER"); user == "" {
	// 	FatalInternal().Msg("Must specify lowpriv user to run as in the env variable USER")
	// }

	if valid, err := util.ValidFile(path); err != nil {
		FatalInternal().Msg("Unable to check " + path + " for validity")
	} else if valid == false {
		FatalInternal().Msg("File " + path + " does not exist or isn't a file")
	}

	run, ok := runner.Runners[lang]
	if !ok {
		FatalInternal().Msg("Language type " + lang + " does not have an assigned runner")
	}

	HandlePrepareError(run.Prepare(path))

	if launchType == LAUNCH_TEST {
		decoder := json.NewDecoder(os.Stdin)
		encoder := json.NewEncoder(os.Stdout)
		for {
			var test string
			if err := decoder.Decode(&test); err == io.EOF {
				break
			} else if err != nil {
				FatalInternal().Msg("Unable to decode test: " + err.Error())
			}
			encoder.Encode(Run(run, test))
		}
	} else if launchType == LAUNCH_SINGLE {
		os.Stdout.WriteString(Run(run, ""))
	} else {
		FatalInternal().Msg("Invalid launch type " + launchType)
	}
}
