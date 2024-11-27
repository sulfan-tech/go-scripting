package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	firestoreDB "go-scripting/configs/firestore"
	"go-scripting/entities"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersRepo interface {
	CreateUser(ctx context.Context, user *entities.User) error
	CreateLogMembership(ctx context.Context, uid string, logMembership entities.UserLogsMembership) error
	QueryUserByPhone(phone string) (*entities.User, error)
	QueryUserByUserAppId(userAppId string) (*entities.User, error)
	GetLatestLogMembership(ctx context.Context, uid string) ([]entities.UserLogsMembership, error)
	GetUserPTPackagesByUID(ctx context.Context, uid string) ([]entities.UserPTPackage, error)
	GetPTPackagesByUID(ctx context.Context, uidMember, uidPackage string) (*entities.UserPTPackage, error)
	GetLogMembershipByExecutionId(ctx context.Context, uid, executionId string) (*entities.UserLogsMembership, error)
	UpdateExpiredMembership(ctx context.Context, uid string, newExpirationDate time.Time) error
	UpdateFieldMembership(ctx context.Context, uid string, updateFields map[string]interface{}) error
	UpdateUserPTPackageFields(ctx context.Context, uid string, packageID string, updateFields map[string]interface{}) error
	UpdateLogMembershipAndTypeChange(ctx context.Context, uid, changeType string) error
}

type FSUser struct {
	client *firestore.Client
}

func NewFirestoreUsersRepository() UsersRepo {
	firestoreClient := firestoreDB.Setup()
	return &FSUser{
		client: firestoreClient,
	}
}

func (r *FSUser) GetPTPackagesByUID(ctx context.Context, uidMember, uidPackage string) (*entities.UserPTPackage, error) {
	var userPTPackages entities.UserPTPackage

	ptPackageDoc := r.client.Collection("users").Doc(uidMember).Collection("vouchers").Doc(uidPackage)

	docSnapshot, err := ptPackageDoc.Get(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := docSnapshot.DataTo(&userPTPackages); err != nil {
		return nil, err
	}

	return &userPTPackages, nil
}

// GetLogMembershipByExecutionId retrieves log membership by ExecutionId for a specific user
func (r *FSUser) GetLogMembershipByExecutionId(ctx context.Context, uid, executionId string) (*entities.UserLogsMembership, error) {
	docRef := r.client.Collection("users").Doc(uid).Collection("logsMembership").Doc(executionId)

	snapshot, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Document not found, continue
		}
		return nil, err
	}

	var logMembership entities.UserLogsMembership
	if err := snapshot.DataTo(&logMembership); err != nil {
		return nil, err
	}

	return &logMembership, nil
}

func (r *FSUser) CreateLogMembership(ctx context.Context, uid string, logMembership entities.UserLogsMembership) error {
	docRef := r.client.Collection("users").Doc(uid).Collection("logsMembership").Doc(logMembership.ExecutionId)

	_, err := docRef.Set(ctx, logMembership)
	if err != nil {
		return err
	}

	return nil
}

func (r *FSUser) QueryUserByUserAppId(userAppId string) (*entities.User, error) {
	ctx := context.Background()

	// Query the "users" collection for the provided phone number
	iter := r.client.Collection("users").Where("userAppId", "==", userAppId).Documents(ctx)

	var user *entities.User

	for {
		// Retrieve the next document (if any)
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			// No documents found
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse the Firestore document into a User struct
		var u entities.User
		if err := doc.DataTo(&u); err != nil {
			return nil, err
		}

		if !u.IsDeleted {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, nil // No valid documents found
	}
	return user, nil
}

func (r *FSUser) QueryUserByPhone(phone string) (*entities.User, error) {
	ctx := context.Background()

	// Query the "users" collection for the provided phone number
	iter := r.client.Collection("users").Where("phone", "==", phone).Documents(ctx)

	var user *entities.User

	for {
		// Retrieve the next document (if any)
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			// No documents found
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse the Firestore document into a User struct
		var u entities.User
		if err := doc.DataTo(&u); err != nil {
			return nil, err
		}

		if !u.IsDeleted {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, nil // No valid documents found
	}
	return user, nil
}

func (r *FSUser) UpdateExpiredMembership(ctx context.Context, uid string, newExpirationDate time.Time) error {
	userRef := r.client.Collection("users").Doc(uid)

	updateData := map[string]interface{}{
		"expiredMembership": newExpirationDate,
	}

	_, err := userRef.Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error updating expiredMembership: %v", err)
	}

	return nil
}

func (r *FSUser) GetLatestLogMembership(ctx context.Context, uid string) ([]entities.UserLogsMembership, error) {
	docRef := r.client.Collection("users").Doc(uid).Collection("logsMembership").
		OrderBy("dateTime", firestore.Desc).
		Limit(1)

	docSnapshot, err := docRef.Documents(ctx).Next()
	if err != nil {
		return nil, err
	}

	var logMemberships []entities.UserLogsMembership
	logMembership := entities.UserLogsMembership{}
	errSnap := docSnapshot.DataTo(&logMembership)
	if errSnap != nil {
		return nil, errSnap
	}

	// Check if dateTime field exists
	if logMembership.DateTime.IsZero() {
		return nil, errors.New("dateTime field does not exist or is empty")
	}

	// Convert dateTime value to time.Time and check if it's valid
	dateTime := logMembership.DateTime
	if dateTime.IsZero() {
		return nil, errors.New("dateTime field is empty")
	}

	logMemberships = append(logMemberships, logMembership)

	return logMemberships, nil
}

func (r *FSUser) UpdateLogMembershipAndTypeChange(ctx context.Context, uid, changeType string) error {
	logMembershipCollectionRef := r.client.Collection("users").Doc(uid).Collection("logsMembership")

	// Query the latest logMembership document based on dateTime
	latestDocSnapshot, err := logMembershipCollectionRef.
		OrderBy("dateTime", firestore.Desc).
		Limit(1).
		Documents(ctx).
		Next()
	if err != nil {
		return err
	}

	// Update the typeChange field of the latest logMembership document
	_, err = latestDocSnapshot.Ref.Update(ctx, []firestore.Update{
		{Path: "typeChange", Value: changeType},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *FSUser) GetUserPTPackagesByUID(ctx context.Context, uid string) ([]entities.UserPTPackage, error) {

	var userPTPackages []entities.UserPTPackage
	ptPackageDocs := r.client.Collection("users").Doc(uid).Collection("vouchers").Where("action", "==", "ptvoucher").Documents(ctx)

	for {
		doc, err := ptPackageDocs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return userPTPackages, err
		}

		ptPackage := entities.UserPTPackage{}
		ptPackage.ID = doc.Ref.ID
		errSnap := doc.DataTo(&ptPackage)
		if errSnap != nil {
			return userPTPackages, errSnap
		}
		userPTPackages = append(userPTPackages, ptPackage)
	}
	return userPTPackages, nil
}

// func (r *FSUser) UpdateUserPTPackage(ctx context.Context, uid string, packageID string, updatedPackage entities.UserPTPackage) error {
// 	docRef := r.client.Collection("users").Doc(uid).Collection("vouchers").Doc(packageID)

// 	_, err := docRef.Set(ctx, updatedPackage, firestore.MergeAll)

// 	return err
// }

func (r *FSUser) CreateUser(ctx context.Context, user *entities.User) error {
	// Construct the reference to the Firestore document
	docRef := r.client.Collection("users").Doc(user.Uid)

	// Set the user data
	_, err := docRef.Set(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

// UpdateUserPTPackageFields updates specific fields of the UserPTPackage document in Firestore.
func (r *FSUser) UpdateUserPTPackageFields(ctx context.Context, uid string, packageID string, updateFields map[string]interface{}) error {
	// Construct the reference to the Firestore document
	docRef := r.client.Collection("users").Doc(uid).Collection("vouchers").Doc(packageID)

	// Update specific fields using Set with MergeAll
	_, err := docRef.Set(ctx, updateFields, firestore.MergeAll)

	return err
}

func (r *FSUser) UpdateFieldMembership(ctx context.Context, uid string, updateFields map[string]interface{}) error {
	docRef := r.client.Collection("users").Doc(uid)

	// Update specific fields using Set with MergeAll
	_, err := docRef.Set(ctx, updateFields, firestore.MergeAll)

	return err
}
