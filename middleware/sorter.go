package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

const SorterKey = "ginson_sort"

func Sorter(ctx *gin.Context) {
	opt := &options.FindOptions{}
	query, exist := ctx.GetQuery("sort")
	if !exist {
		ctx.Set(SorterKey, opt)
		ctx.Next()
		return
	}
	sorts := strings.Split(strings.TrimSpace(query), ",")
	d := bson.D{}
	for _, sort := range sorts {
		if strings.HasPrefix(sort, "-") {
			d = append(d, bson.E{Key: strings.TrimPrefix(sort, "-"), Value: -1})
		} else {
			d = append(d, bson.E{Key: sort, Value: 1})
		}
	}
	opt.SetSort(d)
	ctx.Set(SorterKey, opt)
	ctx.Next()
}

func GetSortOpt(ctx *gin.Context) (sortOpt *options.FindOptions) {
	sorter, exist := ctx.Get(SorterKey)
	if !exist {
		log.Error().Msg("get sort failed, make sure the middleware 'Sorter' is used")
		sortOpt = options.Find()
	} else {
		sortOpt = sorter.(*options.FindOptions)
	}
	return
}
