package links

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

const collection = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
	var l database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	req.ID = primitive.NewObjectID()
	link := database.Link{
		ID:     req.ID,
		URL:    req.URL,
		Title:  req.Title,
		Tags:   req.Tags,
		Images: req.Images,
		UserID: req.UserID,
	}

	_, err := r.db.Collection(collection).InsertOne(ctx, link)
	if err != nil {
		return database.Link{}, err
	}

	return l, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, link, userID string) (database.Link, error) {
	var l database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"url": link, "userid": userID}
	err := r.db.Collection(collection).FindOne(ctx, filter).Decode(&l)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return database.Link{}, nil
		}
		return database.Link{}, err
	}

	return l, nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
	var links []database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{}
	if criteria.UserID != nil {
		filter["userid"] = *criteria.UserID
	}
	if len(criteria.Tags) > 0 {
		filter["tags"] = bson.M{"$all": criteria.Tags}
	}

	findOptions := options.Find()
	if criteria.Limit != nil {
		findOptions.SetLimit(*criteria.Limit)
	}
	if criteria.Offset != nil {
		findOptions.SetSkip(*criteria.Offset)
	}

	cursor, err := r.db.Collection(collection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var l database.Link
		if err := cursor.Decode(&l); err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
