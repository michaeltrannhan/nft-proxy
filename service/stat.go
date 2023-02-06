package services

import (
	nft_proxy "github.com/alphabatem/nft-proxy"
	"github.com/babilu-online/common/context"
)

type StatService struct {
	context.DefaultService

	filesServed    uint64
	requestsServed uint64

	sql *SqliteService
}

const STAT_SVC = "stat_svc"

func (svc StatService) Id() string {
	return STAT_SVC
}

func (svc *StatService) Start() error {
	svc.sql = svc.Service(SQLITE_SVC).(*SqliteService)

	return nil
}

func (svc *StatService) IncrementFileRequests() {
	svc.filesServed++
}

func (svc *StatService) IncrementMediaRequests() {
	svc.requestsServed++
}

func (svc *StatService) ServiceStats() (map[string]interface{}, error) {
	var imgCount int64
	svc.sql.Db().Model(&nft_proxy.SolanaMedia{}).Count(&imgCount)

	return map[string]interface{}{
		"images_stored":   imgCount,
		"files_served":    svc.filesServed,
		"requests_served": svc.requestsServed,
	}, nil
}
