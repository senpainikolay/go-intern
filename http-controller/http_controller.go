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
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write(returnJsonBytesErr("invalid id format or missing"))
		return
	}

	if cachedVal, ok := c.campaignsPerSorceIdCache.Load(sourceId); ok {

		if campaignsJsonBytesCached, ok := cachedVal.([]byte); ok {
			ctx.Write(campaignsJsonBytesCached)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write(returnJsonBytesErr("something wroing with converting the cached value"))
		return
	}

	compagins, err := c.repo.GetCampaignsPerSourceId(sourceId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write(returnJsonBytesErr(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(compagins)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		panic(err)
	}
	c.campaignsPerSorceIdCache.Store(sourceId, jsonBytes)

	ctx.Write(jsonBytes)

}

func returnJsonBytesErr(errStr string) []byte {

	jsonBytes, err := json.Marshal(struct {
		Error bool   `json:"error"`
		Msg   string `json:"msg"`
	}{
		Error: true,
		Msg:   errStr,
	})
	if err != nil {
		panic(err)
	}

	return jsonBytes

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
