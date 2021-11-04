package runner

func CompilerRunner(compiler string) *StepRunner {
	return &StepRunner{
		CompileStep: []string{
			compiler + " -o {{.Out}} {{.File}}",
			"chmod +x {{.Out}}",
		},
		RunStep: "./{{.Out}}",
	}
}

var Runners = map[string]Runner{
	"GCC": CompilerRunner("gcc"),
	"G++": CompilerRunner("g++"),
	"Python": &StepRunner{
		CompileStep: []string{},
		RunStep:     "python {{.File}}",
	},
	"Go": CompilerRunner("go"),
}
