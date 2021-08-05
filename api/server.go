package api

import (
	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/token"
	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		tokenMaker: maker,
		store:      store,
	}

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		val.RegisterValidation("currency", util.CurrencyValidator)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user", server.createNewUserAPI)
	router.POST("/user/login", server.userLoginAPI)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/account", server.createNewAccountAPI)
	authRouter.GET("/account/:id", server.getAccountByIDAPI)
	authRouter.GET("/accounts", server.getListAccountsAPI)
	authRouter.POST("/transfer", server.transferTxAPI)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
