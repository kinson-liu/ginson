package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"strings"
)

const FilterKey = "ginson_filter"

type Query struct {
	Key      string
	Operator string
	Value    string
}

func Filter(ctx *gin.Context) {
	queries := GetQueries(ctx)
	//fmt.Printf("queries: %v\n", queries)

	//for _, query := range queries {
	//	fmt.Printf("%s '%s' %s\n", query.Key, query.Operator, query.Value)
	//}
	filters := ParseQueries(queries)
	//fmt.Printf("filters: %v\n", filters)
	//fmt.Printf("filters length: %v\n", len(filters))

	ctx.Set(FilterKey, filters)
	ctx.Next()
}

func GetQueries(ctx *gin.Context) (queries []Query) {
	urlQuery := ctx.Request.URL.Query()
	innerKeys := []string{"page", "perPage", "sort"}
	sort.Strings(innerKeys)
	//fmt.Printf("urlQuery: %v\n", urlQuery)
	for key, values := range urlQuery {
		index := sort.SearchStrings(innerKeys, key)
		if index < len(innerKeys) && innerKeys[index] == key {
			continue
		}
		if len(values) == 1 {
			query := Query{Key: key}
			index := strings.Index(values[0], "_")
			if index == -1 {
				query.Operator = "eq"
				query.Value = values[0]
			} else {
				query.Operator = values[0][:index]
				query.Value = values[0][index+1:]
			}
			queries = append(queries, query)
		}
	}
	return
}

func ParseQueries(queries []Query) (filters bson.M) {
	filters = bson.M{}
	for _, query := range queries {
		switch query.Operator {
		case "eq", "gt", "gte", "lt", "lte", "ne", "nin":
			filters[query.Key] = bson.M{"$" + query.Operator: query.Value}
		case "like":
			filters[query.Key] = bson.M{"$regex": query.Value}
		case "in":
			filters[query.Key] = bson.M{"$in": strings.Split(query.Value, ",")}
		default:
			filters[query.Key] = bson.M{"$eq": query.Operator + "_" + query.Value}
		}
	}
	return
}

func GetFilterOpt(ctx *gin.Context) (filterOpt bson.M) {
	filter, exist := ctx.Get(FilterKey)
	if !exist {
		log.Error().Msg("get sort failed, make sure the middleware 'Sorter' is used")
		filterOpt = bson.M{}
	} else {
		filterOpt = filter.(bson.M)
	}
	return
}
