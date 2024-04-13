package controller

import (
	"encoding/json"
	"strconv"
	"strings"
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
		return
	}

	if cachedVal, ok := c.campaignsPerSorceIdCache.Load(sourceId); ok {

		if campaignsJsonBytesCached, ok := cachedVal.([]byte); ok {
			ctx.Write(campaignsJsonBytesCached)
			return
		}
	}

	compagins, err := c.repo.GetCampaignsPerSourceId(sourceId)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(compagins)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		panic(err)
	}
	c.campaignsPerSorceIdCache.Store(sourceId, jsonBytes)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(jsonBytes)

}

func (c *GeneralController) GetCampaignsWithDomainsPerSourceIdAndFilterByType(ctx *fasthttp.RequestCtx) {

	sourceIdStr := ctx.QueryArgs().Peek("id")

	ctx.Response.Header.Set("Content-Type", "application/json")

	sourceId, err := strconv.Atoi(string(sourceIdStr))
	if err != nil {
		ctx.Error("invalid id format or missing", fasthttp.StatusBadRequest)
		return
	}

	domainBytes := ctx.QueryArgs().Peek("domain")
	if len(domainBytes) == 0 {
		ctx.Error("no domain specified", fasthttp.StatusInternalServerError)
		return
	}
	domainStr := strings.ToLower(string(domainBytes))
	res := strings.Split(domainStr, ".")
	if len(res) < 2 {
		ctx.Error("invalid domain", fasthttp.StatusInternalServerError)
		return
	}
	domainStr = res[len(res)-2] + "." + res[len(res)-1]

	compagins, err := c.repo.GetCampaignsWithDomainsPerSourceIdAndFilterByType(sourceId, domainStr)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(compagins)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		panic(err)
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(jsonBytes)

}

func Serve(c *GeneralController, port string) {

	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaignsBySource":
			c.GetCampaginsPerSource(ctx)
		case "/capaignsDomainsPerSource":
			c.GetCampaignsWithDomainsPerSourceIdAndFilterByType(ctx)

		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	err := fasthttp.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}

}
