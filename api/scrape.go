package api

import (
	db "Go_web_scrapping/db/sqlc"
	"Go_web_scrapping/token"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//string Response
type stringResponse struct {
	Message string `json:"message"`
}

func stringResFunction(msg string) stringResponse {
	return stringResponse{
		Message: msg,
	}
}

// create scrapping  request
type createScrappingReq struct {
	Url []string `json:"url" binding:"required"`
}

// create scrapping response
type createScrappingRes struct {
	ID              int32     `json:"id"`
	UserId          int32     `json:"user_id"`
	Url             string    `json:"url"`
	ScrappedContent string    `json:"scrapped_content"`
	CreatedAt       time.Time `json:"created_at" `
}

type scrap struct {
	url     string
	content string
}

// Create Scrapping  handles request for Scrapping web
// @Summary Create Scrapping
// @Tags Scrape
// @ID CreateScrapping
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param data body createScrappingReq true "create scrapping request"
// @Success 201 {object} createScrappingRes
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /scrape/create [post]
func (server *Server) createScrapping(ctx *gin.Context) {
	var req createScrappingReq
	var res []createScrappingRes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var wg sync.WaitGroup
	lengthUrl := len(req.Url)
	s := make(chan scrap, lengthUrl)
	wg.Add(lengthUrl)
	for _, u := range req.Url {
		go scrapping(u, s, &wg)
	}
	wg.Wait()
	for range req.Url {
		data := <-s
		arg := db.CreateScrapeParams{
			UserID:   authPayload.UserId,
			Url:      data.url,
			Scrapped: data.content,
		}
		scrappingData, err := server.store.CreateScrape(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		items := createScrappingRes{
			ID:              scrappingData.ID,
			UserId:          scrappingData.UserID,
			Url:             scrappingData.Url,
			ScrappedContent: scrappingData.Scrapped,
			CreatedAt:       scrappingData.CreatedAt,
		}
		res = append(res, items)

	}
	ctx.JSON(http.StatusCreated, res)
}

//web scrapping function
func scrapping(url string, s chan scrap, wg *sync.WaitGroup) {
	defer wg.Done()
	content := url + "_its content static"
	s <- scrap{url, content}

}

//GetScrapedList response
type getScrapedListRes struct {
	ID              int32     `json:"id"`
	UserName        string    `json:"username"`
	Url             string    `json:"url"`
	ScrappedContent string    `json:"scrapped_content"`
	CreatedAt       time.Time `json:"created_at" `
}

func newGetScrapedList(scraped []db.GetScrapeRow) []getScrapedListRes {
	var res []getScrapedListRes
	for _, v := range scraped {
		items := getScrapedListRes{
			ID:              v.ID,
			UserName:        v.Username,
			Url:             v.Url,
			ScrappedContent: v.Scrapped,
			CreatedAt:       v.CreatedAt,
		}
		res = append(res, items)
	}

	return res
}

func mapTotalPage(p int64) []int64 {
	var tPage []int64
	for i := int64(1); i <= p; i++ {
		tPage = append(tPage, i)
	}
	return tPage
}

type getScrapedRes struct {
	Scrape      []getScrapedListRes `json:"Scraped_data"`
	Totalpage   []int64             `json:"total_page"`
	CurrentPage int32               `json:"current_page"`
}

// GetScrapedList handles request for fetch Scraped List
// @Summary get Scraped data of all users
// @Tags Scrape
// @ID getScraped
// @Accept json
// @Produce json
// @Param        page   path      int  true  "page"
// @Success 200 {object} getScrapedRes
// @Success 204 {object} stringResponse
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /list{page} [get]
func (server *Server) getScrapedList(ctx *gin.Context) {

	var res getScrapedRes
	id := ctx.Param("page")
	aid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}

	data, err := server.store.GetScrape(ctx, int32((aid)*10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	count, err := server.store.CountScrape(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(data) == 0 {
		ctx.JSON(http.StatusNoContent, stringResFunction("List is Empty"))
		return
	}
	res = getScrapedRes{
		Scrape:      newGetScrapedList(data),
		Totalpage:   mapTotalPage(10 % count),
		CurrentPage: int32(aid),
	}

	ctx.JSON(http.StatusOK, res)
}

func newGetOwnScrapedList(scraped []db.GetOwnScrapeRow) []getScrapedListRes {
	var res []getScrapedListRes
	for _, v := range scraped {
		items := getScrapedListRes{
			ID:              v.ID,
			UserName:        v.Username,
			Url:             v.Url,
			ScrappedContent: v.Scrapped,
			CreatedAt:       v.CreatedAt,
		}
		res = append(res, items)
	}

	return res
}

// getOwnScrapedList handles request for fetch own Scraped List
// @Summary get Scraped data of  own
// @Tags Scrape
// @ID getownScraped
// @Accept json
// @Produce json
// @Param        page   path      int  true  "page"
// @Security bearerAuth
// @Success 200 {object} getScrapedRes
// @Success 204 {object} stringResponse
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /scrape/{page} [get]
func (server *Server) getOwnScrapedList(ctx *gin.Context) {

	var res getScrapedRes
	id := ctx.Param("page")
	aid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetOwnScrapeParams{
		UserID: authPayload.UserId,
		ID:     int32((aid) * 10),
	}
	data, err := server.store.GetOwnScrape(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(data) == 0 {
		ctx.JSON(http.StatusNoContent, stringResFunction("List is Empty"))
		return
	}
	count, err := server.store.CountScrape(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res = getScrapedRes{
		Scrape:      newGetOwnScrapedList(data),
		Totalpage:   mapTotalPage(10 % count),
		CurrentPage: int32(aid),
	}

	ctx.JSON(http.StatusCreated, res)
}
