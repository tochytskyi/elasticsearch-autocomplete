# elasticsearch-autocomplete

1. `docker-compose up -d` start containers


2. `docker logs elastic01parser -f` follow app logs and wait for indexing success
```
2021/12/02 09:00:19 [200 OK] {
  "name" : "fce1a24664d3",
  "cluster_name" : "docker-cluster",
  "cluster_uuid" : "4rOiayczRz2AcNAOQIpi4g",
  "version" : {
    "number" : "7.15.2",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "93d5a7f6192e8a1a12e154a2b81bf6fa7309da0c",
    "build_date" : "2021-11-04T14:04:42.515624022Z",
    "build_snapshot" : false,
    "lucene_version" : "8.9.0",
    "minimum_wire_compatibility_version" : "6.8.0",
    "minimum_index_compatibility_version" : "6.0.0-beta1"
  },
  "tagline" : "You Know, for Search"
}

2021/12/02 09:00:19 Successfully opened json
2021/12/02 09:00:47 ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔
2021/12/02 09:00:47 Sucessfuly indexed [370,101] documents in 27.504s (13,456 docs/sec)
```


3. `curl --location --request GET 'http://127.0.0.1:9201/words/_count'`
```JSON
{
    "count": 370101,
    "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
    }
}
```


4. Try to search for words like `apple`
```BASH
curl --location --request POST 'http://127.0.0.1:9201/words/_search' \
--header 'Content-Type: application/json' \
--data-raw '{
    "_source": "suggest",
    "suggest": {
        "harry_suggest": {
            "prefix": "appl",
            "completion": {
                "field": "word"
            }
        }
    }
}'
```


Response:
```JSON
{
    "took": 7,
    "timed_out": false,
    "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
    },
    "hits": {
        "total": {
            "value": 0,
            "relation": "eq"
        },
        "max_score": null,
        "hits": []
    },
    "suggest": {
        "harry_suggest": [
            {
                "text": "apple",
                "offset": 0,
                "length": 5,
                "options": [
                    {
                        "text": "apple",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "46588",
                        "_score": 1.0,
                        "_source": {}
                    },
                    {
                        "text": "appleberry",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "363685",
                        "_score": 1.0,
                        "_source": {}
                    },
                    {
                        "text": "appleblossom",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "75793",
                        "_score": 1.0,
                        "_source": {}
                    },
                    {
                        "text": "applecart",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "315619",
                        "_score": 1.0,
                        "_source": {}
                    },
                    {
                        "text": "appled",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "128944",
                        "_score": 1.0,
                        "_source": {}
                    }
                ]
            }
        ]
    }
}
```


5. Try to search for words like `holiday` with 1 mistake
```BASH
curl --location --request POST 'http://127.0.0.1:9201/words/_search' \
--header 'Content-Type: application/json' \
--data-raw '{
    "_source": "suggest",
    "suggest": {
        "harry_suggest": {
            "prefix": "appl",
            "completion": {
                "field": "hoiday"
            }
        }
    }
}'
```

Response:
```JSON
{
    "took": 1,
    "timed_out": false,
    "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
    },
    "hits": {
        "total": {
            "value": 0,
            "relation": "eq"
        },
        "max_score": null,
        "hits": []
    },
    "suggest": {
        "harry_suggest": [
            {
                "text": "hoiday",
                "offset": 0,
                "length": 6,
                "options": []
            }
        ]
    }
}
```


6. Try to search for words like `holiday` with 1 mistake and with `fuzziness` param
```BASH
curl --location --request POST 'http://127.0.0.1:9201/words/_search' \
--header 'Content-Type: application/json' \
--data-raw '{
    "_source": "suggest",
    "suggest": {
        "harry_suggest": {
            "prefix": "appl",
            "completion": {
                "field": "hoiday",
                "fuzzy": {
                    "fuzziness": 1
                }
            }
        }
    }
}'
```

Response:
```JSON
{
    "took": 15,
    "timed_out": false,
    "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
    },
    "hits": {
        "total": {
            "value": 0,
            "relation": "eq"
        },
        "max_score": null,
        "hits": []
    },
    "suggest": {
        "harry_suggest": [
            {
                "text": "hoiday",
                "offset": 0,
                "length": 6,
                "options": [
                    {
                        "text": "holiday",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "326408",
                        "_score": 2.0,
                        "_source": {}
                    },
                    {
                        "text": "holidayed",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "138662",
                        "_score": 2.0,
                        "_source": {}
                    },
                    {
                        "text": "holidayer",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "40663",
                        "_score": 2.0,
                        "_source": {}
                    },
                    {
                        "text": "holidaying",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "144200",
                        "_score": 2.0,
                        "_source": {}
                    },
                    {
                        "text": "holidayism",
                        "_index": "words",
                        "_type": "_doc",
                        "_id": "299286",
                        "_score": 2.0,
                        "_source": {}
                    }
                ]
            }
        ]
    }
}
```