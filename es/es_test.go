package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

type User struct {
	id   string
	name string
	age  int
}

func TestCreateDocsByStruct(t *testing.T) {
	u := User{
		//id:   "200",
		name: "wunder200",
		age:  200,
	}
	_, err := esCli.Index().
		Index("index_name_02"). // 索引名称
		Id(u.name).             // 指定文档id
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

//批量操作
func TestBulkUpdateDocs(t *testing.T) {
	users := []User{
		{
			id:   "1",
			name: "wunder",
			age:  18,
		},
		{
			id:   "2",
			name: "sun",
			age:  20,
		},
	}
	bulkRequest := esCli.Bulk() //初始化新的BulkService。
	for _, u := range users {
		doc := elastic.NewBulkUpdateRequest().Id(u.id).Doc(u).Index("index_name") // 创建一个更新请求
		bulkRequest = bulkRequest.Add(doc)                                        // 添加到批量操作
	}
	_, err := bulkRequest.Do(context.Background())
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
	query = elastic.NewTermQuery("id", "10")
	//
	//// terms
	//query = elastic.NewTermsQuery("field_name", "field_value")
	//
	//match
	//query = elastic.NewMatchQuery("id", "11")
	//
	//// match_phrase
	//query = elastic.NewMatchPhraseQuery("field_name", "field_value")
	//
	//// match_phrase_prefix
	//query = elastic.NewMatchPhrasePrefixQuery("field_name", "field_value")
	//
	////range Gt:大于; Lt:小于; Gte:大于等于; Lte:小于等
	//query = elastic.NewRangeQuery("field_name").Gte(1).Lte(2)
	//
	////regexp
	//query = elastic.NewRegexpQuery("field_name", "regexp_value")

	rets, err := esCli.Search().Index("index_name").Query(query).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//var users []user
	fmt.Printf("Source:%v\n ", string(rets.Hits.Hits[0].Source))

}

// https://blog.csdn.net/Xiao_W_u/article/details/118908282
