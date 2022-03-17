package mongodb

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"stablex/domain"
	"time"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL string, mongoDB string, mongoTimeout int) (domain.OperatorRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) InsertAction(operatorId string, action domain.Action) error {

	id, _ := primitive.ObjectIDFromHex(operatorId)
	update := bson.D{
		{"$push", bson.D{{"actions", bson.D{{"type", action.Type}, {"created_at", action.CreatedAt}}}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("operators")

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)
	if err != nil {
		return errors.Wrap(err, "repository.Operator.Insert")
	}

	fmt.Printf("updated %v doc\n", res.ModifiedCount)
	return nil
}

func (r *mongoRepository) FindOperator(operatorId string, opts string) (*domain.Operator, error) {
	emp := &domain.Operator{}
	return emp, nil
}

func (r *mongoRepository) InsertOperators(ops []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("operators")
	res, err := collection.InsertMany(ctx, ops)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(res)
}
