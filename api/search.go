package api

import (
	db "Go_web_scrapping/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func newSearchList(scraped []db.SearchRow) []getScrapedListRes {
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

// search handles request for search
// @Summary you can search by url
// @Tags Search
// @ID Search
// @Accept json
// @Produce json
// @Success 200 {object} getScrapedListRes
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /search/{q} [get]
func (server *Server) search(ctx *gin.Context) {

	var res []getScrapedListRes
	q := ctx.Param("q")
	data, err := server.store.Search(ctx, q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res = newSearchList(data)

	ctx.JSON(http.StatusCreated, res)
}
