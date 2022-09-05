package es

import (
	"context"
	"github.com/olivere/elastic/v7"
)

var esCli *elastic.Client

func init() {
	var err error
	esCli, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://es:9200"))

	if err != nil {
		panic(err)
	}
}

// 添加
func EsAdd(es_index string, es_body interface{}) string {
	// PS：id是指ES自动生成的id,如果想自定义id可以在Index()后面加.Id(string)。
	// 目前用的是es7.6，所以可以不用加Type默认_doc， 也可以指定方法Type()。
	var id string
	rsAdd, err := esCli.Index().
		Index(es_index).
		BodyJson(es_body).
		Do(context.Background())
	if err != nil {
		panic(err)
	} else {
		id = rsAdd.Id
	}
	return id
}

// EsAddBulk 批量添加
func EsAddBulk(es_index string, es_body []interface{}) {
	bulkRequest := esCli.Bulk()
	for _, v := range es_body {
		tmp := v
		req := elastic.NewBulkIndexRequest().Index(es_index).Doc(tmp)
		bulkRequest = bulkRequest.Add(req)
	}
	_, err := bulkRequest.Do(context.Background())
	if err != nil {
		panic(err)
	}
}

// EsUpdate 修改
func EsUpdate(es_index string, es_id string, es_body interface{}) {
	_, err := esCli.Update().
		Index(es_index).
		Id(es_id).
		Doc(es_body).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
}

// EsUpdateBulk 批量修改
func EsUpdateBulk(es_index string, es_body []map[string]string) {
	// 批量修改和批量增加方法是不同的，批量修改的方法NewBulkUpdateRequest，要调用这个才有DocAsUpsert。
	// 之前在NewBulkIndexRequest下一直没有找到DocAsUpsert。
	//另外用[]map[string]string是因为每个数据要修改的key可能不同,如果要修改的key相同可以用interface 或 struct。
	bulkRequest := esCli.Bulk()
	for _, v := range es_body {
		tmp := v
		tmps := make(map[string]string)
		for is, vs := range tmp {
			if is != "id" {
				tmps[is] = vs
			}
		}
		req := elastic.NewBulkUpdateRequest().Index(es_index).Id(tmp["id"]).Doc(tmps).DocAsUpsert(true)
		bulkRequest = bulkRequest.Add(req)
	}
	_, err := bulkRequest.Do(context.Background())
	if err != nil {
		panic(err)
	}
}

// Delete 删除文档
func Delete(es_index string, id string) {
	_, err := esCli.Delete().Index(es_index).Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
}

// 查询index是否存在
func ExistsIndex(es_index string) bool {
	exists, err := esCli.IndexExists(es_index).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return exists
}

// 查询
func Search(es_index string, page int, li int) interface{} {
	// NewMatchPhraseQuery是代表短语的匹配,Collapse是指定字段相同结果的折叠
	p := (page - 1) * li
	collapsedata := elastic.NewCollapseBuilder("company")
	esq := elastic.NewBoolQuery()
	esq.Must(elastic.NewMatchPhraseQuery("company", "CO LTD."))
	esq.Must(elastic.NewMatchQuery("country", "US"))
	search := esCli.Search().
		Index(es_index).
		From(p).Size(li).
		Query(esq).
		Collapse(collapsedata).
		Pretty(true)
	searchResult, err := search.Do(context.Background())
	if err != nil {
		panic(err)
	} else {
		return searchResult
	}
}
