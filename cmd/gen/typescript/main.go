package main

import (
	"github.com/reading-tribe/anansi/pkg/typescriptgen"
)

func main() {
	generator := typescriptgen.NewGenerator()
	generator.Generate("anansi-api.d.ts")
}
