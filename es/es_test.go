package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Employee struct {
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

//创建索引
func create() {
	//1.使用结构体方式存入到es里面
	e2 := Employee{"jane", "Smith", 20, "I like music", []string{"music"}}
	put, err := esCli.Index().
		Index("test_info").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("indexed %s to index %s, type %s \n", put.Id, put.Index, put.Type)
}

func TestCreateES(t *testing.T) {
	create()
}

func TestCreateDocsByStruct(t *testing.T) {
	//u := User{id: "201", name: "wunder201", age: 201}
	u := User{
		Id:   "202",
		Name: "W202",
		Age:  202,
	}
	_, err := esCli.Index().
		Index("index_name_02"). // 索引名称
		Id(u.Name).             // 指定文档id
		BodyJson(u).            // 	可序列化JSON
		Do(context.Background())

	if err != nil {
		panic(err)
	}
}

func TestCreateDocsByString(t *testing.T) {
	u := `{"name":"wunder1", "age": 11,"id":"10"}`
	_, err := esCli.Index().
		Index("index_name"). // 索引名称
		Id("10").            // 指定文档id
		BodyJson(u).         // 可序列化JSON
		Do(context.Background())

	if err != nil {
		panic(err)
	}
}

func TestUpdateDoc(t *testing.T) {
	_, err := esCli.Update().
		Index("index_name").
		Id("id").
		Doc(map[string]interface{}{"name": "wunder100"}). // 需要修改的字段值
		Do(context.Background())

	if err != nil {
		panic(err)
	}
}

// TestBulkUpdateDocs 批量操作
func TestBulkUpdateDocs(t *testing.T) {
	users := []User{
		{
			Id:   "1000",
			Name: "1000",
			Age:  1000,
		},
		{
			Id:   "2000",
			Name: "2000",
			Age:  2000,
		},
	}
	// 初始化新的BulkService
	bulkService := esCli.Bulk().Index("index_name").Refresh("true")
	for _, u := range users {
		doc := elastic.NewBulkUpdateRequest().
			Id(u.Id).
			Doc(u).Upsert(u)
		bulkService.Add(doc) // 添加到批量操作
	}

	_, err := bulkService.Do(context.Background())
	if err != nil {
		panic(err)
	}

}

func TestQuery(t *testing.T) {
	ret, err := esCli.Get().Index("index_name").
		Id("11"). // 文档id
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("id:%s \n Source:%s ", ret.Id, string(ret.Source))

}

func TestByCondition(t *testing.T) {
	var query elastic.Query

	// match_all
	query = elastic.NewMatchAllQuery()

	// term
	//query = elastic.NewTermQuery("id", "10")
	//
	//// terms
	//query = elastic.NewTermsQuery("field_name", "field_value")
	//
	//match
	query = elastic.NewMatchQuery("IPADDR", "10.4.250.1")
	//
	//// match_phrase
	//query = elastic.NewMatchPhraseQuery("IPADDR", "10.4.250.1")
	//
	//// match_phrase_prefix
	//query = elastic.NewMatchPhrasePrefixQuery("field_name", "field_value")
	//
	////range Gt:大于; Lt:小于; Gte:大于等于; Lte:小于等
	//query = elastic.NewRangeQuery("field_name").Gte(1).Lte(2)
	//
	////regexp
	//query = elastic.NewRegexpQuery("field_name", "regexp_value")

	// 指定返回字段
	rets, err := esCli.Search().
		Index("logstash-oa_login_log*").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("@timestamp", "LOGINTIME", "FROMSOURCE", "USERNAME", "IPADDR")).
		Query(query).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	//var users []user
	fmt.Printf("Source:%v\n ", string(rets.Hits.Hits[0].Source))

}

// https://blog.csdn.net/Xiao_W_u/article/details/118908282
