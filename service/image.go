package services

import (
	"errors"
	"fmt"
	nft_proxy "github.com/alphabatem/nft-proxy"
	"github.com/babilu-online/common/context"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ImageService struct {
	context.DefaultService

	defaultSize int

	httpMedia *http.Client

	solSvc *SolanaImageService
	resize *ResizeService
	sql    *SqliteService
}

const IMG_SVC = "img_svc"

func (svc ImageService) Id() string {
	return IMG_SVC
}

func (svc *ImageService) Start() error {
	svc.solSvc = svc.Service(SOLANA_IMG_SVC).(*SolanaImageService)
	svc.sql = svc.Service(SQLITE_SVC).(*SqliteService)
	svc.resize = svc.Service(RESIZE_SVC).(*ResizeService)

	svc.httpMedia = &http.Client{Timeout: 10 * time.Second}

	svc.defaultSize = 720 //Gifs will be half the size

	return nil
}

func (svc *ImageService) Media(key string) (*nft_proxy.Media, error) {
	if svc.IsSolKey(key) {
		return svc.solSvc.Media(key)
	}

	return nil, errors.New("invalid key")
}

func (svc *ImageService) ImageFile(c *gin.Context, key string) error {
	var media *nft_proxy.Media
	var err error
	if svc.IsSolKey(key) {
		media, err = svc.solSvc.Media(key)
		if err != nil {
			return err
		}
	}

	cacheName := fmt.Sprintf("./cache/solana/%s.%s", media.Mint, media.ImageType)
	ifo, err := os.Stat(cacheName)

	if err != nil || ifo.Size() == 0 { //Missing cached image
		err := svc.fetchMissingImage(media, cacheName)
		if err != nil {
			return err
		}
	}

	log.Printf("Using cached file: %s", cacheName)
	resizedData, err := ioutil.ReadFile(cacheName)
	if err != nil {
		return err
	}

	c.Data(200, fmt.Sprintf("image/%s", media.ImageType), resizedData)
	return nil
}

func (svc *ImageService) fetchMissingImage(media *nft_proxy.Media, cacheName string) error {
	resp, err := svc.httpMedia.Get(media.ImageUri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	output, err := os.Create(cacheName)
	if err != nil {
		return err
	}

	log.Printf("Resizing file: %s", cacheName)
	err = svc.resize.Resize(data, output, svc.defaultSize)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ImageService) MediaFile(c *gin.Context, key string) error {
	var media *nft_proxy.Media
	var err error
	if svc.IsSolKey(key) {
		media, err = svc.solSvc.Media(key)
		if err != nil {
			return err
		}
	}

	if media.MediaUri == "" {
		return errors.New("no media for mint")
	}

	resp, err := svc.httpMedia.Get(media.MediaUri)
	if err != nil {
		return err
	}

	//Write our data
	c.Header("Content-Type", media.MediaType)
	err = resp.Write(c.Writer)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ImageService) IsSolKey(key string) bool {
	_, err := solana.PublicKeyFromBase58(key)
	return err == nil
}
