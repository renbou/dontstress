package runner

func compilerRunner(compiler string) *StepRunner {
	return &StepRunner{
		compileSteps: []string{
			compiler + " -o {{.Out}} {{.File}}",
			"chmod +x {{.Out}}",
		},
		execStep: "./{{.Out}}",
	}
}

func singleStepRunner(cmd string) *StepRunner {
	return &StepRunner{
		compileSteps: []string{},
		execStep:     cmd,
	}
}

var runners = map[string]LangRunner{
	"GCC":    compilerRunner("gcc"),
	"G++":    compilerRunner("g++"),
	"Python": singleStepRunner("python {{.File}}"),
	"Java":   singleStepRunner("java {{.File}}"),
	"Go":     compilerRunner("go"),
}
