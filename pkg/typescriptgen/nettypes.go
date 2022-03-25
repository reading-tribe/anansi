package typescriptgen

import (
	"os"

	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/skia-dev/go2ts"
)

func generateNettypes() {
	f, err := os.Create("nettypes.d.ts")
	if err != nil {
		panic(err)
	}

	generator := go2ts.New()

	// API-Auth
	generator.Add(nettypes.LoginRequest{})
	generator.Add(nettypes.LoginResponse{})
	generator.Add(nettypes.RegisterRequest{})
	generator.Add(nettypes.RegisterResponse{})

	// API-Book
	generator.Add(nettypes.CreateBookRequest{})
	generator.Add(nettypes.CreateBookResponse{})
	generator.Add(nettypes.GetBookResponse{})
	generator.Add(nettypes.ListBooksResponse{})
	generator.Add(nettypes.UpdateBookRequest{})
	generator.Add(nettypes.UpdateBookResponse{})

	// API-Book
	generator.Add(nettypes.CreateTranslationRequest{})
	generator.Add(nettypes.CreateTranslationResponse{})
	generator.Add(nettypes.GetTranslationResponse{})
	generator.Add(nettypes.ListTranslationsResponse{})
	generator.Add(nettypes.UpdateTranslationRequest{})
	generator.Add(nettypes.UpdateTranslationResponse{})

	generator.Render(f)
}
