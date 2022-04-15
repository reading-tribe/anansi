package typescriptgen

import (
	"github.com/reading-tribe/anansi/pkg/nettypes"
)

func (g *Generator) generateNettypes() {
	// API-Auth
	g.generator.Add(nettypes.LoginRequest{})
	g.generator.Add(nettypes.LoginResponse{})
	g.generator.Add(nettypes.RegisterRequest{})
	g.generator.Add(nettypes.RegisterResponse{})
	g.generator.Add(nettypes.LogoutRequest{})
	g.generator.Add(nettypes.RefreshRequest{})

	// API-Book
	g.generator.Add(nettypes.CreateBookRequest{})
	g.generator.Add(nettypes.CreateBookResponse{})
	g.generator.Add(nettypes.GetBookResponse{})
	g.generator.Add(nettypes.ListBooksResponse{})
	g.generator.Add(nettypes.UpdateBookRequest{})
	g.generator.Add(nettypes.UpdateBookResponse{})

	// API-Translation
	g.generator.Add(nettypes.CreateTranslationRequest{})
	g.generator.Add(nettypes.CreateTranslationResponse{})
	g.generator.Add(nettypes.GetTranslationResponse{})
	g.generator.Add(nettypes.ListTranslationsResponse{})
	g.generator.Add(nettypes.UpdateTranslationRequest{})
	g.generator.Add(nettypes.UpdateTranslationResponse{})

	// API-Language
	g.generator.Add(nettypes.ListLanguagesResponse{})

	// API-User
	g.generator.Add(nettypes.ListUsersResponse{})
	g.generator.Add(nettypes.GetUserResponse{})
	g.generator.Add(nettypes.CreateUserResponse{})
	g.generator.Add(nettypes.CreateUserRequest{})
	g.generator.Add(nettypes.UpdateUserRequest{})
	g.generator.Add(nettypes.UpdateUserResponse{})
}
