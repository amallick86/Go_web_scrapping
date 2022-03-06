package api

import (
	db "Go_web_scrapping/db/sqlc"
	"Go_web_scrapping/docs"
	"Go_web_scrapping/token"
	"Go_web_scrapping/util"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	//mode of gin dev mode or release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Go web scrapping"
	docs.SwaggerInfo.Description = "Go_web_scrapping API'S"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	//server route
	server.router = router

	//cors middleware
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		//swager route
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))

		v1.POST("/users", server.createUser)
		v1.POST("/login", server.login)
		v1.GET("/list:page", server.getScrapedList)
		v1.GET("/search:q", server.search)
		v1.POST("/filter", server.filter)

		scrape := v1.Group("/scrape").Use(authMiddleware(server.tokenMaker))
		{
			scrape.POST("/create", server.createScrapping)
			scrape.GET(":page", server.getOwnScrapedList)
		}

	}

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type Err struct {
	Error string `json:"error"`
}

//error response function
func errorResponse(err error) Err {

	return Err{Error: err.Error()}
}
