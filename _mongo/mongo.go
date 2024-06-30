package _mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
)

type MongoServer struct {
	MongoURI string
	DbName   string
	*mongo.Client
	*mongo.Database
	IdLock sync.Mutex
}

type FindCallback func(cur *mongo.Cursor)

func New(mongoURI, database string) (*MongoServer, error) {
	server := MongoServer{MongoURI: mongoURI, DbName: database}
	return server.connect()
}

func (s *MongoServer) connect() (*MongoServer, error) {
	clientOptions := options.Client().ApplyURI(s.MongoURI)
	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	s.Client = client
	s.Database = client.Database(s.DbName)
	//defer client.Disconnect(ctx)
	return s, nil
}

// InsertMany 批量插入Document
func (s *MongoServer) InsertMany(tableName string, documents []interface{}) error {
	// 获取集合引用
	collection := s.Database.Collection(tableName)
	// 插入数据
	_, err := collection.InsertMany(context.Background(), documents)
	return err
}

// InsertOne 插入一条Document
func (s *MongoServer) InsertOne(tableName string, document interface{}) error {
	collection := s.Database.Collection(tableName)
	// 插入数据
	_, err := collection.InsertOne(context.Background(), document)
	return err
}

// InsertOneAndIdV2 插入一条数据，并返回插入的自增ID，注：自增ID需要手动加索引
func (s *MongoServer) InsertOneAndIdV2(tableName string, document interface{}) (int64, error) {
	autoId, err := s.getAutoIncreaseId(tableName)
	if err != nil {
		return 0, err
	}
	// 保存数据
	collection := s.Database.Collection(tableName)
	filter := bson.M{"id": autoId}
	_, err = collection.InsertOne(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	_, err = collection.UpdateOne(context.Background(), filter, bson.M{"$set": document})
	if err != nil {
		return 0, err
	}
	return autoId, err
}

// InsertOneAndId 插入一条数据，并返回插入的自增ID，注：自增ID需要手动加索引，名称：Id
func (s *MongoServer) InsertOneAndId(tableName string, document interface{}) (int64, error) {
	autoId, err := s.getAutoIncreaseId(tableName)
	if err != nil {
		return 0, err
	}
	// 保存数据
	collection := s.Database.Collection(tableName)
	// 自增ID
	err2 := SetAutoId(document, autoId)
	if err2 != nil {
		return 0, err2
	}
	_, err = collection.InsertOne(context.Background(), document)
	return autoId, err
}

// 获取新自增ID
func (s *MongoServer) getAutoIncreaseId(tableName string) (int64, error) {
	s.IdLock.Lock()
	defer s.IdLock.Unlock()
	autoCollection := s.Database.Collection("AutoIncreaseId")
	row, err := autoCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": tableName}, bson.M{"$inc": bson.M{"value": 1}}).DecodeBytes()
	if err != nil && err == mongo.ErrNoDocuments {
		// 插入第一条自增ID
		err = s.InsertOne("AutoIncreaseId", autoIncreaseId{tableName, 1})
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	autoId := row.Lookup("value").Int64() + 1
	return autoId, nil
}

// FindMany 全量查询
func (s *MongoServer) FindMany(tableName string, filter interface{}, opts ...*options.FindOptions) []bson.Raw {
	collection := s.Database.Collection(tableName)
	// 查询数据
	cursor, err := collection.Find(context.Background(), filter, opts...)
	if err != nil {
		return nil
	}
	// 转为切片
	var raws []bson.Raw
	for cursor.Next(context.Background()) {
		raws = append(raws, cursor.Current)
	}
	return raws
}

// FindManyResult 全量查询，自己解析数据
func (s *MongoServer) FindManyResult(tableName string, filter interface{}, callback FindCallback, opts ...*options.FindOptions) error {
	collection := s.Database.Collection(tableName)
	// 查询数据
	cursor, err := collection.Find(context.Background(), filter, opts...)
	if err != nil {
		return err
	}
	// 转为切片
	for cursor.Next(context.Background()) {
		callback(cursor)
	}
	return nil
}

// FindOne 查询一条
func (s *MongoServer) FindOne(tableName string, filter interface{}) bson.Raw {
	collection := s.Database.Collection(tableName)
	// 查询数据
	res := collection.FindOne(context.Background(), filter)
	bytes, err := res.DecodeBytes()
	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return nil, nil
		//} else {
		return nil
		//}
	}
	return bytes
}

// FindOneResult 查询一条，自己解析数据
func (s *MongoServer) FindOneResult(tableName string, filter interface{}) *mongo.SingleResult {
	collection := s.Database.Collection(tableName)
	// 查询数据
	return collection.FindOne(context.Background(), filter)
}

// UpdateOne 只更新一条
func (s *MongoServer) UpdateOne(tableName string, update interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	collection := s.Database.Collection(tableName)
	return collection.UpdateOne(context.Background(), filter, update)
}

// UpdateMany 全量更新
func (s *MongoServer) UpdateMany(tableName string, update interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	collection := s.Database.Collection(tableName)
	return collection.UpdateMany(context.Background(), filter, update)
}

// Count 统计条数
func (s *MongoServer) Count(tableName string, filter interface{}) (int64, error) {
	collection := s.Database.Collection(tableName)
	return collection.CountDocuments(context.Background(), filter)
}

// DeleteOne 删除一条数据
func (s *MongoServer) DeleteOne(tableName string, filter interface{}) error {
	collection := s.Database.Collection(tableName)
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

// DeleteMany 全量删除数据
func (s *MongoServer) DeleteMany(tableName string, filter interface{}) error {
	collection := s.Database.Collection(tableName)
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}

// Aggregate 聚合查询
func (s *MongoServer) Aggregate(tableName string, pipeline interface{}, opts ...*options.AggregateOptions) []bson.Raw {
	collection := s.Database.Collection(tableName)
	// 查询数据
	cursor, err := collection.Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		return nil
	}
	// 转为切片
	var raws []bson.Raw
	for cursor.Next(context.Background()) {
		raws = append(raws, cursor.Current)
	}
	return raws
}

// SetAutoId 插入自增ID
func SetAutoId(document interface{}, id int64) error {
	value := reflect.ValueOf(document)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("document must be a pointer to a struct")
	}
	idField := value.Elem().FieldByName("Id")
	if !idField.IsValid() || !idField.CanSet() {
		return fmt.Errorf("id field is not accessible: <Id>")
	}
	switch idField.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		idField.SetInt(id)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		idField.SetUint(uint64(id))
	default:
		return fmt.Errorf("ID field has unsupported type")
	}
	return nil
}

func (s *MongoServer) FindOneAndUpdate(tableName string, update interface{}, filter interface{}, opts ...*options.FindOneAndUpdateOptions) bson.Raw {
	collection := s.Database.Collection(tableName)
	// 查询数据
	res := collection.FindOneAndUpdate(context.Background(), filter, update, opts...)
	bytes, err := res.DecodeBytes()
	if err != nil {
		return nil
	}
	return bytes
}
