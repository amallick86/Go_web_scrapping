package api

import (
	db "Go_web_scrapping/db/sqlc"
	"Go_web_scrapping/token"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

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

// Create Scrapping  handles request for user creation
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
	wg.Add(1)
	go scrapping(req.Url, s, &wg)
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
func scrapping(url []string, s chan scrap, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, u := range url {
		content := u + "_its content static"
		s <- scrap{u, content}
	}

}

// id := ctx.Param("id")
// 	aid, err := strconv.Atoi(id)
// 	if err != nil {
// 		httputil.NewError(ctx, http.StatusBadRequest, err)
// 		return
// 	}
