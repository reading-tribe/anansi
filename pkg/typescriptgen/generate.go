package typescriptgen

import (
	"os"

	"github.com/skia-dev/go2ts"
)

type Generator struct {
	generator *go2ts.Go2TS
}

func NewGenerator() *Generator {
	generator := go2ts.New()

	return &Generator{
		generator: generator,
	}
}

func (g *Generator) Generate(outputfilename string) {
	f, err := os.Create(outputfilename)
	if err != nil {
		panic(err)
	}

	g.generateNettypes()
	g.generateIdx()

	g.generator.Render(f)

	g.generator = go2ts.New()
}
