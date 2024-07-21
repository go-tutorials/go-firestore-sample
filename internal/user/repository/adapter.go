package repository

import (
	"context"
	"reflect"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	fs "github.com/core-go/firestore"

	"go-service/internal/user/model"
)

type UserAdapter struct {
	Collection *firestore.CollectionRef
	Map        map[string]string
}

func NewUserAdapter(client *firestore.Client) *UserAdapter {
	userType := reflect.TypeOf(model.User{})
	maps := fs.MakeFirestoreMap(userType)
	collection := client.Collection("users")
	return &UserAdapter{Collection: collection, Map: maps}
}

func (a *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	iter := a.Collection.Documents(ctx)
	var users []model.User
	for {
		doc, er1 := iter.Next()
		if er1 == iterator.Done {
			break
		}
		if er1 != nil {
			return nil, er1
		}
		var user model.User
		er2 := doc.DataTo(&user)
		if er2 != nil {
			return users, er2
		}

		user.Id = doc.Ref.ID
		user.CreateTime = &doc.CreateTime
		user.UpdateTime = &doc.UpdateTime
		users = append(users, user)
	}
	return users, nil
}

func (a *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	doc, er1 := a.Collection.Doc(id).Get(ctx)
	var user model.User
	if er1 != nil {
		return nil, er1
	}
	er2 := doc.DataTo(&user)
	if er2 == nil {
		user.Id = id
		user.CreateTime = &doc.CreateTime
		user.UpdateTime = &doc.UpdateTime
	}
	return &user, er2
}

func (a *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	var docRef *firestore.DocumentRef
	if len(user.Id) > 0 {
		docRef = a.Collection.Doc(user.Id)
	} else {
		docRef = a.Collection.NewDoc()
		user.Id = docRef.ID
	}
	res, err := docRef.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "Document already exists") {
			return 0, nil
		} else {
			return 0, err
		}
	}
	user.CreateTime = &res.UpdateTime
	user.UpdateTime = &res.UpdateTime
	return 1, nil
}

func (a *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	docRef := a.Collection.Doc(user.Id)
	doc, er1 := docRef.Get(ctx)
	if er1 != nil {
		if strings.HasSuffix(er1.Error(), " not found") {
			return 0, nil
		}
		return -1, er1
	}
	res, er2 := docRef.Set(ctx, user)
	if er2 != nil {
		return -1, er2
	}
	user.CreateTime = &doc.CreateTime
	user.UpdateTime = &res.UpdateTime
	return 1, nil
}

func (a *UserAdapter) Patch(ctx context.Context, json map[string]interface{}) (int64, error) {
	uid := json["id"]
	id := uid.(string)
	docRef := a.Collection.Doc(id)
	doc, er1 := docRef.Get(ctx)
	if er1 != nil {
		return -1, er1
	}
	delete(json, "id")

	dest := fs.MapToFirestore(json, doc.Data(), a.Map)
	_, er2 := docRef.Set(ctx, dest)
	if er2 != nil {
		return -1, er2
	}
	return 1, nil
}

func (a *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	_, err := a.Collection.Doc(id).Delete(ctx)
	if err != nil {
		return -1, err
	}
	return 1, nil
}
