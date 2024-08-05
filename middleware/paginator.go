package middleware

import (
	"fmt"
	"ginson/api/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PaginatorKey = "ginson_paginator"

func Paginator(ctx *gin.Context) {
	log.Debug().Msg("get paginator")
	page := Page{}
	err := ctx.ShouldBindQuery(&page)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	ctx.Set(PaginatorKey, page)
	ctx.Next()
}

type Page struct {
	Page    int64 `form:"page"`
	PerPage int64 `form:"perPage"`
}

func (p *Page) Offset() int64 {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 10
	}
	return (p.Page - 1) * p.PerPage
}

func (p *Page) Options() (opt *options.FindOptions) {
	opt = &options.FindOptions{}
	opt.SetSkip(p.Offset())
	opt.SetLimit(p.PerPage)
	return
}

func GetPageOpt(ctx *gin.Context) (pageOpt *options.FindOptions) {
	page := Page{
		Page:    1,
		PerPage: 10,
	}
	paginator, exist := ctx.Get(PaginatorKey)
	if !exist {
		log.Error().Msg("get page failed, make sure the middleware 'Paginator' is used")
	} else {
		page = paginator.(Page)
	}
	pageOpt = page.Options()
	fmt.Printf("page: %v\n", page)
	return
}

func GetPageFromOptions(opts ...*options.FindOptions) (page Page) {
	for _, opt := range opts {
		fmt.Printf("opt: %+v\n", opt)
		if opt == nil {
			continue
		}
		if opt.Limit != nil {
			page.PerPage = *opt.Limit
			fmt.Printf("PerPage: %+v\n", page.PerPage)
		}
		if opt.Skip != nil {
			if page.PerPage == 0 {
				continue
			} else {
				page.Page = *opt.Skip/page.PerPage + 1
				fmt.Printf("Page: %+v\n", page.Page)

			}
		}
	}
	return
}
