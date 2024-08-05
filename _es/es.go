package _es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"sync"
)

// Smartcn 分词器模式
const (
	Smartcn = "smartcn"
)

// IK分词器模式
const (
	IKSmart   = "ik_smart"    //用得多
	IkMaxWord = "ik_max_word" //IK细分词
)

// ES默认分词器 无法分词中文
const Standard = "standard"

type EsServer struct {
	EsURI    string
	UserName string
	PassWord string
	*elastic.Client
	//*mongo.Database
	IdLock sync.Mutex
}

type FindData struct {
	Id  string //唯一值，利用她进行增删改查
	Raw []byte //对应的值
}

type failedDetail struct {
	Ids         []string
	ErrorMagmap map[string]string
}

func NewEsServer(mongoURI, userName string, passWord string) (*EsServer, error) {
	server := EsServer{EsURI: mongoURI, UserName: userName, PassWord: passWord}
	return server.connect()
}

func (s *EsServer) connect() (*EsServer, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(s.EsURI), elastic.SetBasicAuth(s.UserName, s.PassWord), elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	s.Client = client
	return s, nil
}

// 定义映射
func (s *EsServer) CreateMapping(indexName string, mapping string) error {
	_, err := s.Client.CreateIndex(indexName).BodyString(mapping).Do(context.Background())
	return err
}

// 更新映射， 仅支持添加字段, 已有字段无法修改
func (s *EsServer) AddMapping(indexName string, mapping string) error {
	_, err := s.Client.PutMapping().Index(indexName).BodyString(mapping).Do(context.Background())
	return err
}

// 插入一条数据 无需自定义id可为空
func (s *EsServer) PutIndexOne(indexName string, id string, documents interface{}) error {
	// 保存文档到 Elasticsearch
	body := s.Client.Index().
		Index(indexName).
		Refresh("true").
		BodyJson(documents)
	if id != "" {
		body = body.Id(id)
	}
	_, err := body.Do(context.Background())
	return err
}

func (s *EsServer) PutIndexMany(indexName string, documents map[string]interface{}) (success []string, fDetail map[string]string, err error) {
	// 保存文档到 Elasticsearch
	req := s.Client.Bulk().Index(indexName)
	for id, item := range documents {
		if id != "" {
			fmt.Printf("正在导入ES数据索引:index:%s Id:%s \n", indexName, id)
			doc := elastic.NewBulkIndexRequest().Id(id).Doc(item)
			req.Add(doc)
		}
	}

	bulkResponse, err := req.Refresh("true").Do(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if bulkResponse == nil {
		return nil, nil, fmt.Errorf("expected bulkResponse to be != nil; got nil")
	}

	if req.NumberOfActions() != 0 {
		return nil, nil, fmt.Errorf("CreateBulkUsers expected bulkRequest.NumberOfActions %d; got %d", 0, req.NumberOfActions())
	}

	success = make([]string, 0)       //成功的ID
	fDetail = make(map[string]string) //失败的ID和失败信息

	for _, v := range bulkResponse.Items {
		var item = *v["index"]
		if item.Error == nil {
			success = append(success, item.Id)
		} else {
			b, _ := json.Marshal(item.Error)
			fDetail[item.Id] = string(b)
		}
	}

	return success, fDetail, err
}

// UpdateUserById ... 根据id新增或更新数据(单条) 仅更新传入的字段
func (s *EsServer) UpdateUserById(indexName string, id string, update interface{}) (err error) {
	_, err = s.Client.Update().
		Index(indexName).
		Id(id).
		Refresh("true").
		// update为结构体或map, 需注意的是如果使用结构体零值也会去更新原记录
		Upsert(update).
		// true 无则插入, 有则更新, 设置为false时记录不存在将报错
		DocAsUpsert(true).
		Do(context.Background())
	return err
}

// DeleteUserById 指定id删除数据  DELETE /user/_doc/2
func (s *EsServer) DeleteUserById(indexName string, id string) (err error) {
	deleteResult, err := s.Client.Delete().Index(indexName).Id(id).Refresh("true").Do(context.Background())
	if err != nil {
		return err
	}
	if deleteResult == nil {
		return errors.New(fmt.Sprintf("expected result to be != nil; got: %v", deleteResult))
	}

	// 检查文档是否存在
	exists := s.RowExists(indexName, id)
	if exists {
		return errors.New(fmt.Sprintf("expected exists %v; got %v \n", false, exists))
	}
	return nil
}

// RowExists 判断id 的数据是否存在
func (s *EsServer) RowExists(indexName string, id string) bool {
	exists, err := s.Client.Exists().Index(indexName).Id(id).Do(context.TODO())
	if err != nil {
		fmt.Printf("err---->:%v \n", err)
		return false
	}
	return exists
}

// skip:第几页，pageSize 返回几条，几条 elastic.NewMatchQuery(fieldName, text)
func (s *EsServer) WorkSegmentQuery(indexName string, filter *elastic.MatchQuery, page int, pageSize int) ([]FindData, error) {
	searchResult, err := s.Client.Search(indexName).Query(filter).From((page - 1) * pageSize).Size(pageSize).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if searchResult.Hits.TotalHits.Value > 0 {
		var raws []FindData
		for _, hit := range searchResult.Hits.Hits {
			raws = append(raws, FindData{hit.Id, hit.Source})
		}
		return raws, nil
	}
	// 处理没有匹配结果的情况
	return nil, fmt.Errorf("没有查询到数据")
}

// 对文本进行分解 analyzer分解器
func (s *EsServer) IndexAnalyze(analyzer string, text string) ([]string, error) {
	service, err := s.Client.IndexAnalyze().Analyzer(analyzer).Text(text).Do(context.Background())
	if err != nil {
		return nil, err
	}
	var tokens []string
	for _, token := range service.Tokens {
		tokens = append(tokens, token.Token)
	}
	return tokens, nil
}

// GetUserInfo ...  GET users/_doc/1 获取指定id的数据
func (s *EsServer) GetInfoOne(indexName string, id string) ([]byte, error) {
	get, err := s.Client.Get().Index(indexName).Id(id).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !get.Found {
		return nil, fmt.Errorf("没有找到")
	}
	//转为json格式
	source, err := get.Source.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return source, err
}

//
