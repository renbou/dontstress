package compile

func compilerRunner(cmd, ext string) *StepCompiler {
	return &StepCompiler{
		steps: []string{
			cmd + " -o {{.Out}} {{.File}}",
		},
		extension: ext,
	}
}

var compilers = map[string]LangCompiler{
	"GCC": compilerRunner("gcc", ".c"),
	"G++": compilerRunner("g++", ".cpp"),
	"Go":  compilerRunner("go", ".go"),
}

// Gets compiler for lang. If compiler is not defined returns false as second value
func GetCompiler(lang string) (LangCompiler, bool) {
	compiler, ok := compilers[lang]
	return compiler, ok
}
