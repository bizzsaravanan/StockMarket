package db

import (
	"context"
	"errors"
	"reflect"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	DatabaseName string
}
type SC mongo.SessionContext
type Func func(SC) error
type Collection = *mongo.Collection
type M = primitive.M
type D = primitive.D
type ObjectId = primitive.ObjectID

var Client *mongo.Client
var C map[string]*mongo.Collection

func Connect() error {
	Ctx := context.Background()
	var err error
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(viper.GetString("MongoUrl")))
	if err != nil {
		return err
	}
	C = make(map[string]*mongo.Collection)
	Client = client
	return nil
}

func (m *MongoRepo) Col(e interface{}) *mongo.Collection {
	cname := typeName(e)
	if m.DatabaseName == "" {
		m.DatabaseName = viper.GetString("MongoDBName")
	}
	db := Client.Database(m.DatabaseName)
	r2 := db.Collection(cname)
	return r2
}

func typeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	if isSlice(t) {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return t.Name()
}

func isSlice(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice
}

func isPtr(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr
}

func (m *MongoRepo) CreateIndex(e interface{}, field string, unique bool) error {
	ctx := context.Background()
	defer ctx.Done()

	cl := m.Col(e)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},    // 1 for ascending order
		Options: options.Index().SetUnique(unique), // Set unique index option
	}

	_, err := cl.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (m *MongoRepo) CreateCompoundIndex(e interface{}, fields ...string) error {
	ctx := context.Background()
	defer ctx.Done()

	cl := m.Col(e)

	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1}) // 1 for ascending order, use -1 for descending
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true), // Set unique option if needed
	}

	_, err := cl.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (m *MongoRepo) InsertOne(e interface{}) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(e)
	_, err := cl.InsertOne(ctx, e)
	return err
}

func (m *MongoRepo) InsertMany(e []interface{}, c interface{}) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(c)
	_, err := cl.InsertMany(ctx, e)
	return err
}

func (m *MongoRepo) FindOne(i interface{}, q bson.M) error {
	if !isPtr(i) {
		return errors.New("must pass pointer")
	}
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	err := cl.FindOne(ctx, q).Decode(i)
	return err
}

func (m *MongoRepo) FindAll(i interface{}, q bson.M, sort bson.M) error {
	if !isPtr(i) {
		return errors.New("must pass pointer")
	}
	opts := options.Find()
	opts.SetSort(sort)
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	cursor, err := cl.Find(ctx, q, opts)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, i)
	return err
}

func (m *MongoRepo) FindAllWithProjection(i interface{}, q bson.M, sort bson.M, p bson.M) error {
	if !isPtr(i) {
		return errors.New("must pass pointer")
	}
	opts := options.Find()
	opts.SetSort(sort)
	opts.SetProjection(p)
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	cursor, err := cl.Find(ctx, q, opts)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, i)
	return err
}

func (m *MongoRepo) FindAllPagination(i interface{}, q bson.M, sort bson.M, page int64, size int64) error {
	if !isPtr(i) {
		return errors.New("must pass pointer")
	}
	opts := options.Find()
	var skip int64 = 0
	if page > 1 {
		skip = size * (page - 1)
	}
	opts.SetSort(sort)
	opts.SetLimit(size)
	opts.SetSkip(skip)
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	cursor, err := cl.Find(ctx, q, opts)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, i)
	return err
}

func (m *MongoRepo) FindAllPaginationWithProjection(i interface{}, q bson.M, sort bson.M, p bson.M, page int64, size int64) error {
	if !isPtr(i) {
		return errors.New("must pass pointer")
	}
	opts := options.Find()
	var skip int64 = 0
	if page > 1 {
		skip = size * (page - 1)
	}
	opts.SetSort(sort)
	opts.SetProjection(p)
	opts.SetLimit(size)
	opts.SetSkip(skip)
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	cursor, err := cl.Find(ctx, q, opts)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, i)
	return err
}

func (m *MongoRepo) FindAndUpdate(i interface{}, q bson.M, update bson.M, opts *options.FindOneAndUpdateOptions) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	err := cl.FindOneAndUpdate(ctx, q, update, opts).Decode(i)
	return err
}

func (m *MongoRepo) UpdateOne(i interface{}, q bson.M, update bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	_, err := cl.UpdateOne(ctx, q, update)
	return err
}

func (m *MongoRepo) UpdateMany(i interface{}, q bson.M, update bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	_, err := cl.UpdateMany(ctx, q, update)
	return err
}

func (m *MongoRepo) DeleteOne(i interface{}, q bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	_, err := cl.DeleteOne(ctx, q)
	return err
}

func (m *MongoRepo) DeleteMany(i interface{}, q bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := m.Col(i)
	_, err := cl.DeleteMany(ctx, q)
	return err
}
