package main

import (
	"errors"
	"log"

	nft_proxy "github.com/alphabatem/nft-proxy"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
)

type collectionLoader struct {
	metaWorkerCount  int
	fileWorkerCount  int
	mediaWorkerCount int

	done chan struct{} // Add channel for graceful shutdown

	metaDataIn chan *token_metadata.Metadata
	fileDataIn chan *nft_proxy.NFTMetadataSimple
	mediaIn    chan *nft_proxy.Media
}

func main() {
	log.Printf("Loading collection images: %s", "TODO")

	l := collectionLoader{
		metaWorkerCount:  3,
		fileWorkerCount:  3,
		mediaWorkerCount: 1,
		metaDataIn:       make(chan *token_metadata.Metadata),
		fileDataIn:       make(chan *nft_proxy.NFTMetadataSimple),
		mediaIn:          make(chan *nft_proxy.Media),
	}

	l.spawnWorkers()

	//TODO Get collection
	err := l.loadCollection()
	if err != nil {
		panic(err)
	}

	//TODO Fetch all the mints for that collection
	//TODO Fetch Mints/Hash List

	//TODO Batch into batches of 100
	//TODO Pass to metaDataIn<-

	//TODO Fetch all the metadata accounts for that collection
	//TODO Fetch all images for the accounts
	//TODO Fetch Image
	//TODO Resize Image 500x500
	//TODO Fetch Media
}

func (l *collectionLoader) spawnWorkers() {
	spawnWorker := func(worker func(), count int) {
		for i := 0; i < count; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in worker: %v", r)
					}
				}()
				worker()
			}()
		}
	}
	spawnWorker(l.metaDataWorker, l.metaWorkerCount)
	spawnWorker(l.fileDataWorker, l.fileWorkerCount)
	spawnWorker(l.mediaWorker, l.mediaWorkerCount)
}

func (l *collectionLoader) loadCollection() error {
	return nil
}

// Fetches the off-chain data from the on-chain account & passes to `fileDataWorker`
func (l *collectionLoader) metaDataWorker() {
	return
}

// Downloads required files & passes to `mediaWorker`
func (l *collectionLoader) fileDataWorker() {
	return
}

// Stores media data down to SQL with proper error handling
func (l *collectionLoader) mediaWorker() {
	for {
		select {
		case m := <-l.mediaIn:
			if err := l.saveMedia(m); err != nil {
				log.Printf("Failed to save media: %v", err)
			}
		case <-l.done:
			return
		}
	}
}

func (l *collectionLoader) saveMedia(m *nft_proxy.Media) error {
	if m == nil {
		return errors.New("nil media object")
	}
	log.Printf("Saving media: %s", m.MediaUri)
	// TODO: Implement actual DB save
	return nil
}
