package runner

var Runners = map[string]Runner{
	"GCC-11": &StepRunner{
		CompileStep: []string{
			"gcc -o {{.Out}} {{.File}}",
			"chmod +x {{.Out}}",
		},
		RunStep: "./{{.Out}}",
	},
	"G++-11": &StepRunner{
		CompileStep: []string{
			"g++ -o {{.Out}} {{.File}}",
			"chmod +x {{.Out}}",
		},
		RunStep: "./{{.Out}}",
	},
	"Python3.9": &StepRunner{
		CompileStep: []string{},
		RunStep:     "python {{.File}}",
	},
	"Go1.17.2": &StepRunner{
		CompileStep: []string{
			"go build -o {{.Out}} {{.File}}",
			"chmod +x {{.Out}}",
		},
		RunStep: "./{{.Out}}",
	},
}
