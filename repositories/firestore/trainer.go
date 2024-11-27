package firestore

import (
	"context"
	"fmt"

	firestoreDB "go-scripting/configs/firestore"
	"go-scripting/entities"

	"cloud.google.com/go/firestore"
)

type PersonalTrainerRepo struct {
	client *firestore.Client
}

type PersonalTrainer interface {
	GetPT(ctx context.Context, email string) (*entities.PersonalTrainer, error)
}

func NewFirestoreTrainerRepository() PersonalTrainer {
	firestoreClient := firestoreDB.Setup()
	return &PersonalTrainerRepo{
		client: firestoreClient,
	}
}

func (t *PersonalTrainerRepo) GetPT(ctx context.Context, email string) (*entities.PersonalTrainer, error) {
	query := t.client.Collection("trainers").Where("email", "==", email).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query Firestore: %w", err)
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("personal trainer not found")
	}

	var trainer entities.PersonalTrainer
	err = docs[0].DataTo(&trainer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Firestore document: %w", err)
	}

	return &trainer, nil
}
