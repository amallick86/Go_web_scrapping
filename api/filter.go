package api

import (
	db "Go_web_scrapping/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// filter request
type filterReq struct {
	FromDate string `json:"from_date" binding:"required"`
	ToDate   string `json:"to_date" `
}

func newfilterList(scraped []db.FilterRow) []getScrapedListRes {
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

// filter handles request for filter
// @Summary you can filter by time
// @Tags Filter
// @ID filter
// @Accept json
// @Produce json
// @Param data body filterReq true "filter request"
// @Success 200 {object} getScrapedListRes
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /filter [post]
func (server *Server) filter(ctx *gin.Context) {
	var req filterReq
	var res []getScrapedListRes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}
	if req.FromDate == "" {
		minDate, err := server.store.MinDate(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.FromDate = minDate.Format("2006-01-02")
	}
	if req.ToDate == "" {
		now := time.Now()
		req.ToDate = now.Format("2006-01-02")
	}
	layout := "2006-01-02"
	from, err := time.Parse(layout, req.FromDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}

	to, err := time.Parse(layout, req.ToDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}

	arg := db.FilterParams{
		CreatedAt:   from,
		CreatedAt_2: to,
	}

	data, err := server.store.Filter(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res = newfilterList(data)

	ctx.JSON(http.StatusCreated, res)
}
