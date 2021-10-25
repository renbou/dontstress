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
	"GCC-11": CompilerRunner("gcc"),
	"G++-11": CompilerRunner("g++"),
	"Python3.9": &StepRunner{
		CompileStep: []string{},
		RunStep:     "python {{.File}}",
	},
	"Go1.17.2": CompilerRunner("go"),
}
