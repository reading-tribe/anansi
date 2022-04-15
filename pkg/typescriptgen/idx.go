package typescriptgen

import (
	"github.com/reading-tribe/anansi/pkg/idx"
)

func (g *Generator) generateIdx() {
	g.generator.Add(idx.BookID(""))
	g.generator.Add(idx.ChildProfileID(""))
	g.generator.Add(idx.DiversityAndInclusionID(""))
	g.generator.Add(idx.PageID(""))
	g.generator.Add(idx.SessionID(""))
	g.generator.Add(idx.TranslationID(""))
	g.generator.Add(idx.UserID(""))
}
