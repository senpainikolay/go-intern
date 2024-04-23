package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/senpainikolay/go-tasks/models"
	"github.com/senpainikolay/go-tasks/repository"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
	_ "github.com/valyala/fasthttp/pprofhandler"
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

	var wg sync.WaitGroup
	wg.Add(len(compagins.Campaigns))

	priceChan := make(chan models.MinPriceCampaign)
	done := make(chan models.MinPriceCampaign)

	go func() {
		minPrice := models.MinPriceCampaign{Price: 9999999}
		for p := range priceChan {
			if p.Price < minPrice.Price {
				minPrice = p
			}
		}
		done <- minPrice
	}()

	for _, val := range compagins.Campaigns {
		go func(cName string) {
			defer wg.Done()
			minPrice := callSleep()
			priceChan <- models.MinPriceCampaign{Name: cName, Price: minPrice}
		}(val.Name)
	}

	wg.Wait()
	close(priceChan)

	minPriceCampaign := <-done

	finalRes := models.CampaignsWithSelectedMinPriceCampaign{
		Campaigns:        compagins,
		MinPriceCampaign: minPriceCampaign,
	}

	jsonBytes, err := json.Marshal(finalRes)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		panic(err)
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(jsonBytes)

}

func callSleep() int {
	timeToSleep := rand.Intn(10)
	randomPrice := rand.Intn(1000) + 100
	time.Sleep(time.Duration(timeToSleep) * time.Second)
	fmt.Printf("Sleeped: %v \n", timeToSleep)
	return randomPrice
}

func Serve(c *GeneralController, port string) {

	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaignsBySource":
			c.GetCampaginsPerSource(ctx)
		case "/capaignsDomainsPerSource":
			c.GetCampaignsWithDomainsPerSourceIdAndFilterByType(ctx)
		case "/debug/pprof/profile":
			pprofhandler.PprofHandler(ctx)
		case "/debug/pprof/heap":
			pprofhandler.PprofHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	err := fasthttp.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}

}
