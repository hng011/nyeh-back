package firebase

import (
	"context"
	"nyeh-back/internal/core"
	"nyeh-back/internal/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type firestoreContact struct {
	Email     string `firestore:"email"`
	Github    string `firestore:"github"`
	Linkedin  string `firestore:"linkedin"`
	Instagram string `firestore:"instagram"`
}

type firestoreBio struct {
	Name    string           `firestore:"name"`
	AboutMe string           `firestore:"about_me"`
	Address string           `firestore:"address"`
	Contact firestoreContact `firestore:"contact"`
}

type firestoreBioRepo struct {
	client *firestore.Client
}

func NewFirestoreRepo(client *firestore.Client) domain.BioRepository {
	return &firestoreBioRepo{client: client}
}

func (f *firestoreBioRepo) GetBio(ctx context.Context, email string) (*domain.Bio, error) {
	doc, err := f.client.Collection(core.Settings.FIREBASE_FS_COLL_BIO).Doc(email).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var fBio firestoreBio
	if err := doc.DataTo(&fBio); err != nil {
		return nil, err
	}

	return &domain.Bio{
		Name:    fBio.Name,
		AboutMe: fBio.AboutMe,
		Address: fBio.Address,
		Contact: domain.Contact{
			Email:     fBio.Contact.Email,
			Github:    fBio.Contact.Github,
			Linkedin:  fBio.Contact.Linkedin,
			Instagram: fBio.Contact.Instagram,
		},
	}, nil
}

func (f *firestoreBioRepo) UpsertBio(ctx context.Context, email string, bio *domain.Bio) error {
	// init value and pass the reference to the collection set
	fBio := firestoreBio{
		Name:    bio.Name,
		AboutMe: bio.AboutMe,
		Address: bio.Address,
		Contact: firestoreContact{
			Email:     bio.Contact.Email,
			Github:    bio.Contact.Github,
			Linkedin:  bio.Contact.Linkedin,
			Instagram: bio.Contact.Instagram,
		},
	}

	_, err := f.client.Collection(core.Settings.FIREBASE_FS_COLL_BIO).Doc(email).Set(ctx, fBio)
	return err
}

func (f *firestoreBioRepo) DeleteBio(ctx context.Context, email string) error {
	_, err := f.client.Collection(core.Settings.FIREBASE_FS_COLL_BIO).Doc(email).Delete(ctx)
	return err
}
