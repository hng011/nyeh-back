package domain

import "context"

type Contact struct {
	Email     string
	Github    string
	Linkedin  string
	Instagram string
}

type Bio struct {
	Name    string
	AboutMe string
	Contact Contact
	Address string
}

type BioRepository interface {
	GetBio(ctx context.Context, email string) (*Bio, error)
	UpsertBio(ctx context.Context, email string, bio *Bio) error
	DeleteBio(ctx context.Context, email string) error
}

type BioUsecase interface {
	Get(ctx context.Context, email string) (*Bio, error)
	Upsert(ctx context.Context, email string, bio *Bio) error
	Delete(ctx context.Context, email string) error
}
