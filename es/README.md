# 0 连接ES
```shell
options := []elastic.ClientOptionFunc{
	elastic.SetURL("http://xxxxxxx:9200"),
	elastic.SetSniff(true),      //是否开启集群嗅探
	elastic.SetHealthcheckInterval(10 * time.Second), //设置两次运行状况检查之间的间隔, 默认60s
	elastic.SetGzip(false),  //启用或禁用gzip压缩
	elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),  //ERROR日志输出配置
	elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),  //INFO级别日志输出配置
}
options = append(options, elastic.SetBasicAuth(
	"xxxx",            //账号
	"xxxxxxxxxxxxxx",  //密码
))
con, err := elastic.NewClient(options...)

```
# 1 索引操作
## 1.1 创建索引
```shell
type mi = map[string]interface{}
mapping := mi{
	"settings": mi{
		"number_of_shards":   3,
		"number_of_replicas": 2,
	},
	"mappings": mi{
		"_doc": mi{  //type名
			"properties": mi{
				"id": mi{  //整形字段, 允许精确匹配
					"type": "integer",
				},
				"name": mi{
					"type":            "text",  //字符串类型且进行分词, 允许模糊匹配
					"analyzer":        "ik_smart", //设置分词工具
					"search_analyzer": "ik_smart",
					"fields": mi{    //当需要对模糊匹配的字符串也允许进行精确匹配时假如此配置
						"keyword": mi{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"date_field": mi{  //时间类型, 允许精确匹配
					"type": "date",
				},
				"keyword_field": mi{ //字符串类型, 允许精确匹配
					"type": "keyword",
				},
				"nested_field": mi{ //嵌套类型
					"type": "nested",
					"properties": mi{
						"id": mi{
							"type": "integer",
						},
						"start_time": mi{ //长整型, 允许精确匹配
							"type": "long",
						},
						"end_time": mi{
							"type": "long",
						},
					},
				},
			},
		},
	},
}
indexName := "xxxxxxxxxxxxxxxxxxxxx"  //要创建的索引名
_, err = conf.ES().CreateIndex(indexName).BodyJson(mapping).Do(context.Background())

```
## 1.2 判断索引是否存在
```shell
//exists true 表示索引已存在
exists, err := conf.ES().IndexExists(indexName).Do(context.Background())
```

## 1.3 更新索引
 仅支持添加字段, 已有字段无法修改

```shell
type mi = map[string]interface{}
mapping := mi{ 
	"properties": mi{
		"id": mi{  //整形字段, 允许精确匹配
			"type": "integer",
		},
	},
}
_, err = conf.ES().PutMapping().
	Index(indexName).
	Type("_doc").
	BodyJson(mapping).
	Do(context.Background())

```

## 1.4 删除索引
```shell
_, err = conf.ES().DeleteIndex(indexName).Do(context.Background())
```

## 1.5 删除迁移
将一个索引的数据迁移到另一个索引中, 一般用于索引结构发生改变时使用新索引存储数据
```shell
type mi = map[string]interface{}
_, err = conf.ES().Reindex().Body(mi{
	"source": mi{
		"index": oldIndexName,
	},
	"dest": mi{
		"index": newIndexName,
	},
}).Do(context.Background())

```

## 1.6 设置index别名
设置index别名后就可以使用别名查询数据
```shell
_, err = conf.ES().Alias().Action(
	elastic.NewAliasAddAction(oldIndexName).Index(newIndexName),
).Do(context.Background())
```

# 2 数据操作
## 2.1  新增或覆盖数据(单条)
此操作相同id的数据会被覆盖
```shell
_, err = conf.ES().Index().
Index(indexName).
Type("_doc").
// id为字符串, 创建一条此id的数据或覆盖已有此id的记录
// data为结构体或map, 当然结构需要跟索引的mapping类型保持一致
Id(id).BodyJson(data).
Do(context.Background())

```

## 2.2 根据id新增或更新数据(单条)
仅更新传入的字段, 而不是像 4.1 进行整条记录覆盖
```shell
_, err = conf.ES().Update().
	Index(t.index()).
	Type("_doc").
	Id(id).
	// data为结构体或map, 需注意的是如果使用结构体零值也会去更新原记录
	Upsert(data).
	// true 无则插入, 有则更新, 设置为false时记录不存在将报错
	DocAsUpsert(true). 
	Do(context.Background())

```

## 2.3 根据id新增或更新数据(批量)
```shell
bulkRequest := conf.ES().Bulk()
// data map[int]interface{}, key为id, value为要更新的数据
for id, v := range data {
	doc := elastic.NewBulkUpdateRequest().
		Index(t.index()).
		Type("_doc").
		Id(strconv.Itoa(id)).
		Doc(v).
		// true 无则插入, 有则更新, 设置为false时记录不存在将报错
		DocAsUpsert(true)
	bulkRequest.Add(doc)
}
bulkResponse, err := bulkRequest.Do(context.Background())
if err != nil {
	return
}
// 获取操作失败的记录
bad := bulkResponse.Failed()
if len(bad) > 0 {
	s, _ := jsoniter.MarshalToString(bad)
	err = errors.New("部分记录更新失败 " + s)
}

```

## 2.4 根据条件更新数据
```shell
_, err = conf.ES().UpdateByQuery().
	Index(indexName).
	Type("_doc").
	//查询条件, 详细配置查询条件请查看章节 5
	Query(query).
	//要执行的更新操作, 详细配置请查看章节 6及7.1
	Script(script).
	Do(context.Background()) 
```

## 2.5 查询
```shell
_, err = conf.ES().Search().
	Index(indexName).
	//偏移量
	From(0).
	//返回数据的条数
	Size(10).
	//指定返回数据的字段(此处指定返回id和name), 全部返回则无需设置
	FetchSourceContext(elastic.NewFetchSourceContext(true).Include("id", "name")).
	//查询条件, 详细配置查询条件请查看章节 5
	Query(query).
	//按照id升序排序, 无需排序则可跳过此设置, 多个Sort会按先后顺序依次生效
	Sort("id", true).
	//自定义排序规则, 详细写法请查看章节 6及7.2
	SortBy(sorter).
	Do(context.Background())

```

## 3 查询条件Query设置
## 3.1 一个示例
```shell
{
    "bool": {
        "filter": [
            {
                "nested": {
                    "path": "nested_field",
                    "query": {
                        "range": {
                            "nested_field.start_time": {
                                "from": 1581475200,
                                "include_lower": true,
                                "include_upper": true,
                                "to": null
                            }
                        }
                    }
                }
            },
            {
                "nested": {
                    "path": "nested_field",
                    "query": {
                        "range": {
                            "nested_field.end_time": {
                                "from": null,
                                "include_lower": true,
                                "include_upper": true,
                                "to": 1581481440
                            }
                        }
                    }
                }
            }
        ],
        "must": {
            "terms": {
                "id": [
                    4181,
                    4175
                ]
            }
        }
    }
}

```
实现上述查询条件的go代码如下
```shell
query := elastic.NewBoolQuery()
query.Must(elastic.NewTermsQuery("id", []int{4181, 4175}))
query.Filter(elastic.NewNestedQuery(
	"nested_field",
	// nested_field.start_time >= 1581475200
	elastic.NewRangeQuery("nested_field.start_time").Gte(1581475200),
))
query.Filter(elastic.NewNestedQuery(
	"nested_field",
	// nested_field.start_time <= 1581481440
	elastic.NewRangeQuery("nested_field.end_time").Lte(1581481440),
))

```

## 3.2 match 模糊匹配
```shell
// name字段模糊匹配
elastic.NewMatchQuery("name", val)
```

## 3.3 terms 精确匹配
```shell
// name字段精确匹配
elastic.NewTermsQuery("name.keyword", val...)
```

## 3.4 range 范围匹配
```shell
// id >= 10, id <= 100
elastic.NewRangeQuery("id").Gte(10).Lte(100)
```

## 3.5 nested 嵌套结构查询
```shell
elastic.NewNestedQuery(
	"nested_field",
	query,  //此处query中的字段 都需要加上nested_field前缀, 比如 nested_field.id
)
```

## 4 常用的查询
## 4.1 时间查询）
筛选starttime字段,查找大于2020/05/13 18:38:21,并且小于2020/05/14 18:38:21的数据
```shell
type Task struct {
	TaskID    string `json:"taskid"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Count     int    `json:"count"`
}

var typ Task
boolQ := elastic.NewBoolQuery()
//生成查询语句,筛选starttime字段,查找大于2020/05/13 18:38:21,并且小于2020/05/14 18:38:21的数据
boolQ.Filter(elastic.NewRangeQuery("starttime").Gte("2020/05/13 18:38:21"), 
              elastic.NewRangeQuery("starttime").Lte("2020/05/14 18:38:21")
            )
res, _ := Es.Client.Search("task").Type("doc").Query(boolQ).Do(context.Background())
    //从搜索结果中取数据的方法
    for _, item := range res.Each(reflect.TypeOf(typ)) {
        if t, ok := item.(Task); ok {
            fmt.Println(t)
        }
    }
```

## 4.2 查询包含关键字的查询方法
```shell
type Task struct {
	TaskID    string `json:"taskid"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Count     int    `json:"count"`
}

//查找包含"cifs"的所有数据
func (Es *Elastic) FindKeyword() {
    //因为不确定cifs如何出现，可能是cifs01，也可能是01cifs，所以采用这种方法
    keyword := "cifs"
    keys := fmt.Sprintf("name:*%s*", keyword)
    boolQ.Filter(elastic.NewQueryStringQuery(keys))
     res, _ := Es.Client.Search("task").Type("doc").Query(boolQ).Do(context.Background())
        //从搜索结果中取数据的方法
        for _, item := range res.Each(reflect.TypeOf(typ)) {
            if t, ok := item.(Task); ok {
                fmt.Println(t)
            }
        }
}
```

## 4.3 多条件查询
如果说我们现在不仅仅需要找到符合时间的，也需要找到符合关键字的查询
```shell
func (Es *Elastic) FindAll() {
    //因为不确定cifs如何出现,可能是cifs01,也可能是01cifs,所以采用这种方法
    boolQ := elastic.NewBoolQuery()
    keyword := "cifs"
    keys := fmt.Sprintf("name:*%s*", keyword)
    boolQ.Filter(elastic.NewRangeQuery("starttime").Gte("2020/05/13 18:38:21"), 
                 elastic.NewRangeQuery("starttime").Lte("2020/05/14 18:38:21"), 
                 elastic.NewQueryStringQuery(keys)
                 )
	res, err := Es.Client.Search("task").Type("doc").Query(boolQ).Do(context.Background())
        //从搜索结果中取数据的方法
        for _, item := range res.Each(reflect.TypeOf(typ)) {
            if t, ok := item.(Task); ok {
                fmt.Println(t)
            }
        }
}
```

## 4.4 统计
```shell
func (Es *Elastic) GetTaskLogCount() (int, error) {
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewRangeQuery("starttime").Gte("2020/05/13 18:38:21"), 
	             elastic.NewRangeQuery("starttime").Lte("2020/05/14 18:38:21"))
	//统计count
	count, err := Es.Client.Count("task").Type("doc").Query(boolQ).Do(context.Background())
	if err != nil {
		return 0, nil
	}
	return int(count), nil
}
```

