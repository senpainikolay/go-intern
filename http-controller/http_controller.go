package controller

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/senpainikolay/go-tasks/repository"
	"github.com/valyala/fasthttp"
)

type GeneralController struct {
	repo                     *repository.GeneralRepository
	campaignsPerSorceIdCache sync.Map
}

func NewController(repo *repository.GeneralRepository) *GeneralController {
	return &GeneralController{
		repo:                     repo,
		campaignsPerSorceIdCache: sync.Map{},
	}
}

func (c *GeneralController) GetCampaginsPerSource(ctx *fasthttp.RequestCtx) {

	sourceIdStr := ctx.QueryArgs().Peek("id")

	ctx.Response.Header.Set("Content-Type", "application/json")

	sourceId, err := strconv.Atoi(string(sourceIdStr))
	if err != nil {
		ctx.Error("invalid id format or missing", fasthttp.StatusBadRequest)
	}

	if cachedVal, ok := c.campaignsPerSorceIdCache.Load(sourceId); ok {

		if campaignsJsonBytesCached, ok := cachedVal.([]byte); ok {
			ctx.Write(campaignsJsonBytesCached)
			return
		}
		ctx.Error("something wroing with converting the cached value", fasthttp.StatusInternalServerError)
	}

	compagins, err := c.repo.GetCampaignsPerSourceId(sourceId)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	jsonBytes, err := json.Marshal(compagins)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		panic(err)
	}
	c.campaignsPerSorceIdCache.Store(sourceId, jsonBytes)

	ctx.Write(jsonBytes)

}

func Serve(c *GeneralController, port string) {

	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaignsBySource":
			c.GetCampaginsPerSource(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	err := fasthttp.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}

}
