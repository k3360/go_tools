package _mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type MongoServer struct {
	MongoURI string
	DbName   string
	*mongo.Client
	*mongo.Database
	IdLock sync.Mutex
}

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

// 批量插入Document
func (s *MongoServer) InsertMany(tableName string, documents []interface{}) error {
	// 获取集合引用
	collection := s.Database.Collection(tableName)
	// 插入数据
	_, err := collection.InsertMany(context.Background(), documents)
	return err
}

// 插入一条Document
func (s *MongoServer) InsertOne(tableName string, document interface{}) error {
	collection := s.Database.Collection(tableName)
	// 插入数据
	_, err := collection.InsertOne(context.Background(), document)
	return err
}

// 插入一条数据，并返回插入的自增ID
func (s *MongoServer) InsertOneAndId(tableName string, document interface{}) (int64, error) {
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

// 全量查询
func (s *MongoServer) FindMany(tableName string, filter bson.M) ([]bson.Raw, error) {
	collection := s.Database.Collection(tableName)
	// 查询数据
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	// 转为切片
	var raws []bson.Raw
	for cursor.Next(context.Background()) {
		raws = append(raws, cursor.Current)
	}
	return raws, nil
}

// 查询一条
func (s *MongoServer) FindOne(tableName string, filter bson.M) (bson.Raw, error) {
	collection := s.Database.Collection(tableName)
	// 查询数据
	res := collection.FindOne(context.Background(), filter)
	bytes, err := res.DecodeBytes()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return bytes, err
}

// 只更新一条
func (s *MongoServer) UpdateOne(tableName string, update bson.M, filter bson.M) error {
	collection := s.Database.Collection(tableName)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

// 全量更新
func (s *MongoServer) UpdateMany(tableName string, update bson.M, filter bson.M) error {
	collection := s.Database.Collection(tableName)
	_, err := collection.UpdateMany(context.Background(), filter, update)
	return err
}

// 统计条数
func (s *MongoServer) Count(tableName string, filter bson.M) (int64, error) {
	collection := s.Database.Collection(tableName)
	return collection.CountDocuments(context.Background(), filter)
}