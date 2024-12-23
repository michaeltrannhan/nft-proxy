package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/babilu-online/common/context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title NFT Aggregator Swap API
// @version 1.1.27
// @description NFT Aggregator Swap API
// @schemes https

// @host https://agg.alphabatem.com
// @BasePath /
// @query.collection.format multi

type HttpService struct {
	context.DefaultService
	BaseURL string
	Port    int

	imgSvc  *ImageService
	statSvc *StatService

	defaultImage []byte
}

var ErrUnauthorized = errors.New("unauthorized")
var DeleteResponseOK = `{"status": 200, "error": ""}`

func (svc HttpService) Id() string {
	return "http"
}

func (svc *HttpService) Configure(ctx *context.Context) error {
	port := os.Getenv("HTTP_PORT")
	portFlag, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	svc.Port = portFlag

	svc.defaultImage, err = ioutil.ReadFile("./docs/failed_image.jpg")
	if err != nil {
		return err
	}

	return svc.DefaultService.Configure(ctx)
}

func (svc *HttpService) Start() error {
	svc.imgSvc = svc.Service(IMG_SVC).(*ImageService)
	svc.statSvc = svc.Service(STAT_SVC).(*StatService)

	r := gin.Default()

	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	//r.Static("static", "static")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//Validation endpoints
	r.GET("/ping", svc.ping)
	r.GET("/stats", svc.stats)

	v1 := r.Group("/v1")
	//docs.SwaggerInfo.BasePath = "/v1"

	// unify repeated routes for tokens and nfts
	svc.registerNFTEndpoints(v1, "tokens")
	svc.registerNFTEndpoints(v1, "nfts")

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r.Run(fmt.Sprintf(":%v", svc.Port))
}

func (svc *HttpService) registerNFTEndpoints(g *gin.RouterGroup, prefix string) {
	r := g.Group(prefix)
	r.GET("/:id", svc.showNFT)
	r.GET("/:id/image", svc.showNFTImage) // So much repetition but same service
	r.GET("/:id/media", svc.showNFTMedia)
}

type Pong struct {
	Message string `json:"message"`
}

// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router /ping [get]
func (svc *HttpService) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// @Summary Ping liquify service
// @Accept  json
// @Produce json
// @Router /stats [get]
func (svc *HttpService) stats(c *gin.Context) {
	stats, err := svc.statSvc.ServiceStats()
	if err != nil {
		svc.paramErr(c, err)
		return
	}

	c.JSON(200, stats)
}

// @Summary Get NFT metadata
// @Description Get NFT metadata by ID
// @Accept  json
// @Produce json
// @Param   id  path  string  true  "NFT ID"
// @Router /v1/nfts/{id} [get]
func (svc *HttpService) showNFT(c *gin.Context) {
	svc.statSvc.IncrementMediaRequests()

	skipCache, _ := strconv.ParseBool(c.DefaultQuery("nocache", ""))
	if skipCache || rand.Intn(1000) == 1 {
		if err := svc.imgSvc.ClearCache(c.Param("id")); err != nil {
			svc.paramErr(c, err)
			return
		}
	}

	media, err := svc.imgSvc.Media(c.Param("id"), skipCache)
	if err != nil {
		svc.paramErr(c, err)
		return
	}

	c.Header("Cache-Control", "public, max-age=172800")
	c.Header("Expires", time.Now().AddDate(0, 0, 2).Format(http.TimeFormat))

	c.JSON(200, media)
}

// @Summary Get NFT image
// @Description Get NFT image by ID
// @Accept  json
// @Produce image/*
// @Param   id  path  string  true  "NFT ID"
// @Router /v1/nfts/{id}/image [get]
func (svc *HttpService) showNFTImage(c *gin.Context) {
	svc.statSvc.IncrementImageFileRequests()
	err := svc.imgSvc.ImageFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}

// @Summary Get NFT media
// @Description Get NFT media file by ID
// @Accept  json
// @Produce */*
// @Param   id  path  string  true  "NFT ID"
// @Router /v1/nfts/{id}/media [get]
func (svc *HttpService) showNFTMedia(c *gin.Context) {
	svc.statSvc.IncrementMediaFileRequests()
	err := svc.imgSvc.MediaFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}

// Consistent error handling with proper status codes
func (svc *HttpService) paramErr(c *gin.Context, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, ErrUnauthorized) {
		status = http.StatusUnauthorized
	}
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

// TODO Replace with placeholder image
func (svc *HttpService) mediaError(c *gin.Context, err error) {
	log.Printf("Media Error: %v", err)
	c.Header("Cache-Control", "public, max-age=60, must-revalidate")
	c.Data(http.StatusOK, "image/jpeg", svc.defaultImage)
}
