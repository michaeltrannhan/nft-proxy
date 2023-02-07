package services

import (
	"encoding/json"
	nft_proxy "github.com/alphabatem/nft-proxy"
	"github.com/babilu-online/common/context"
	"github.com/gagliardetto/solana-go"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type SolanaImageService struct {
	context.DefaultService
	sql *SqliteService
	sol *SolanaService

	http *http.Client
}

const SOLANA_IMG_SVC = "solana_img_svc"

func (svc SolanaImageService) Id() string {
	return SOLANA_IMG_SVC
}

func (svc *SolanaImageService) Start() error {
	svc.http = &http.Client{Timeout: 5 * time.Second}

	svc.sql = svc.Service(SQLITE_SVC).(*SqliteService)
	svc.sol = svc.Service(SOLANA_SVC).(*SolanaService)
	return nil
}

func (svc *SolanaImageService) Media(key string) (*nft_proxy.Media, error) {
	var media *nft_proxy.SolanaMedia
	err := svc.sql.Db().First(&media, "mint = ?", key).Error
	if err != nil {
		log.Printf("Fetching metadata for: %s - %s", key, err)
		media, err = svc.FetchMetadata(key)
		if err != nil {
			return nil, err //Still cant get metadata
		}
	}

	return media.Media(), nil
}

func (svc *SolanaImageService) FetchMetadata(key string) (*nft_proxy.SolanaMedia, error) {
	metadata, err := svc.retrieve(key)
	if err != nil {
		return nil, err
	}

	media, err := svc.cache(key, metadata, "")
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (svc *SolanaImageService) retrieve(key string) (*nft_proxy.NFTMetadataSimple, error) {
	pk, err := solana.PublicKeyFromBase58(key)
	if err != nil {
		return nil, err
	}
	tokenData, err := svc.sol.TokenData(pk)
	if err != nil {
		log.Printf("No token data")
		return nil, err
	}

	return svc.retrieveFile(tokenData.Data.Uri)
}

func (svc *SolanaImageService) retrieveFile(uri string) (*nft_proxy.NFTMetadataSimple, error) {
	file, err := svc.http.Get(strings.Trim(uri, "\x00")) //Strip crap off urls
	if err != nil {
		return nil, err
	}

	if file.StatusCode != 200 {
		return nil, err
	}

	defer file.Body.Close()
	data, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return nil, err
	}

	var metadata nft_proxy.NFTMetadataSimple
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (svc *SolanaImageService) cache(key string, metadata *nft_proxy.NFTMetadataSimple, localPath string) (*nft_proxy.SolanaMedia, error) {
	media := nft_proxy.SolanaMedia{
		Mint:      key,
		LocalPath: localPath,
	}

	if metadata != nil {
		media.ImageUri = metadata.Image
		media.ImageType = svc.guessImageType(metadata)

		mediaFile := metadata.AnimationFile()
		if mediaFile != nil {
			media.MediaUri = mediaFile.URL
			mediaFile.Type = "mp4"
			if strings.Contains(mediaFile.Type, "/") {
				media.MediaType = strings.Split(mediaFile.Type, "/")[1]
			}
		}
	}

	return &media, svc.sql.Db().Clauses(clause.OnConflict{DoNothing: true}).Create(&media).Error
}

func (svc *SolanaImageService) guessImageType(metadata *nft_proxy.NFTMetadataSimple) string {
	if metadata == nil {
		return "jpg"
	}

	imageType := ""
	imgFile := metadata.ImageFile()
	if imgFile != nil && strings.Contains(imgFile.Type, "/") {
		imageType = strings.Split(imgFile.Type, "/")[1]
	}
	if imageType == "" {
		parts := strings.Split(metadata.Image, ".")
		lastPart := parts[len(parts)-1]
		if strings.Contains(lastPart, "=") {
			parts := strings.Split(lastPart, "=")
			imageType = parts[len(parts)-1]
		} else {
			imageType = lastPart
		}
	}

	if !svc.ValidType(imageType) {
		log.Printf("Invalid image type guessed: %s", imageType)
		return "jpg"
	}

	return imageType
}

func (svc *SolanaImageService) ValidType(imageType string) bool {
	switch imageType {
	case "png":
		return true
	case "jpg":
		return true
	case "jpeg":
		return true
	case "gif":
		return true
	case "svg":
		return true
	default:
		return false
	}
}
