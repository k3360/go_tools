package mongodb

type autoIncreaseId struct {
	Key   string `bson:"_id"`
	Value int64  `bson:"value"`
}
