package database

import (
	"context"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userDatabase struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) domain.UserRepository {
	return &userDatabase{
		client: client,
	}
}

func (db *userDatabase) userCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(db.client, "user")
	return userCollection
}

func (db *userDatabase) Create(ctx context.Context, user domain.User) (string, error) {
	userDTO := mapUserToUserDTO(user)

	_, err := db.userCollection(ctx).InsertOne(ctx, userDTO)
	if err != nil {
		return "", err
	}
	return user.UserID, nil
}

func (db *userDatabase) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	var userDTO UserDTO

	err := db.userCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", userID}}).Decode(&userDTO)
	if err != nil {
		return nil, err
	}

	user := mapUserDTOToUser(userDTO)

	return &user, nil
}

func (db *userDatabase) Search(ctx context.Context, filters domain.UserFilters) ([]domain.User, error) {

	filter := bson.M{}
	if len(filters.Role) > 0 {
		filter["role"] = filters.Role
	}
	if len(filters.Email) > 0 {
		filter["email"] = filters.Email
	}
	if len(filters.Username) > 0 {
		filter["username"] = filters.Username
	}
	if len(filters.FirstName) > 0 {
		filter["first_name"] = filters.FirstName
	}
	if *filters.Active {
		filter["active"] = true
	}
	if len(filters.UserID) > 0 {
		filter["_id"] = filters.UserID
	}

	var usersResult UserDTO
	err := db.userCollection(ctx).FindOne(ctx, filter).Decode(&usersResult)

	if err != nil {
		return nil, err
	}

	users := mapUserDTOsToUsers([]UserDTO{usersResult})

	return users, nil
}

func (db *userDatabase) Update(ctx context.Context, userToUpdate domain.User) error {
	userDTO := mapUserToUserDTO(userToUpdate)

	filter := bson.M{
		"_id": userToUpdate.UserID,
	}

	result := db.userCollection(ctx).FindOneAndUpdate(ctx, filter, userDTO)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (db *userDatabase) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.NewValidationError("userID cannot be empty", nil)
	}

	_, err := db.userCollection(ctx).DeleteOne(ctx, bson.D{{"_id", userID}})
	if err != nil {
		return err
	}

	return nil
}
