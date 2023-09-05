package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mhg14/simplebank/db/sqlc"
	"github.com/mhg14/simplebank/token"
	"github.com/mhg14/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) setupRouter() {
	router := gin.Default()

	// add routes to router
	router.POST("/users", s.createUserHandler)
	router.POST("/users/login", s.loginUserHandler)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccountHandler)
	authRoutes.GET("/accounts/:id", s.getAccountHandler)
	authRoutes.GET("/accounts", s.listAccountsHandler)
	authRoutes.PATCH("/accounts/:id", s.updateAccountHandler)
	authRoutes.DELETE("/accounts/:id", s.deleteAccountHandler)

	authRoutes.POST("/transfers", s.createTransferHandler)

	s.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
