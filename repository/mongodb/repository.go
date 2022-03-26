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
	"log"
	"stablex/domain"
	"stablex/domain/helper"
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

// InsertAction - insert action
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

// FindOperator - find operator
func (r *mongoRepository) FindOperator(operatorId string, opts domain.OperatorFilter) (*domain.Operator, error) {
	var operator *domain.Operator

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(operatorId)

	collection := r.client.Database(r.database).Collection("operators")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&operator)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return operator, nil
}

func (r *mongoRepository) GetOperators(opts domain.OperatorFilter) ([]*domain.Operator, error) {
	var from time.Time
	var to time.Time
	var dateFilter = bson.M{}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("operators")

	if opts.Current {
		from = helper.GetUTCDate(time.Now(), 0)
		dateFilter["$gte"] = from

		to = helper.GetUTCDate(from, 24)
		dateFilter["$gte"] = to

	} else {
		if !opts.FromDate.IsZero() {
			from = helper.GetUTCDate(opts.FromDate, 0)
			dateFilter["$gte"] = from
		}
		if !opts.ToDate.IsZero() {
			to = helper.GetUTCDate(opts.ToDate, 0)
			dateFilter["$gte"] = to
		}
	}

	filters := bson.M{"actions": bson.M{"$elemMatch": bson.M{"created_at": dateFilter}}}

	cursor, err := collection.Find(ctx, filters)
	if err != nil {
		log.Fatal(err)
	}

	var ops []*domain.Operator
	for cursor.Next(ctx) {
		var opt *domain.Operator
		if err := cursor.Decode(&opt); err != nil {
			log.Fatal(err)
		}
		ops = append(ops, opt)
	}
	return ops, nil
}

// InsertAction - insert action
func (r *mongoRepository) UpdateOperator(operatorId string, opts domain.OperatorFilter) error {

	id, _ := primitive.ObjectIDFromHex(operatorId)

	fields := bson.M{}
	if opts.Password != "" {
		fields["password"] = opts.Password
	}
	if opts.Password != "" {
		fields["position"] = opts.Position
	}

	update := bson.M{"$set": fields}

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
	fmt.Printf("----res/: %v", res)
	return nil
}

// InsertOperators - seeder function
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
