package usecase

import (
	"context"
	"nyeh-back/internal/domain"
)

type bioUsecase struct {
	bioRepo domain.BioRepository
}

func NewBioUsecase(repo domain.BioRepository) domain.BioUsecase {
	return &bioUsecase{bioRepo: repo}
}

func (u *bioUsecase) Get(ctx context.Context, email string) (*domain.Bio, error) {
	return u.bioRepo.GetBio(ctx, email)
}

func (u *bioUsecase) Upsert(ctx context.Context, email string, bio *domain.Bio) error {
	return u.bioRepo.UpsertBio(ctx, email, bio)
}

func (u *bioUsecase) Delete(ctx context.Context, email string) error {
	return u.bioRepo.DeleteBio(ctx, email)
}
