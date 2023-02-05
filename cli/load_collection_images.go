package main

import (
	nft_proxy "github.com/alphabatem/nft-proxy"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"log"
)

type collectionLoader struct {
	metaWorkerCount  int
	fileWorkerCount  int
	mediaWorkerCount int

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
	for i := 0; i < l.metaWorkerCount; i++ {
		go l.metaDataWorker()
	}
	for i := 0; i < l.fileWorkerCount; i++ {
		go l.fileDataWorker()
	}
	for i := 0; i < l.mediaWorkerCount; i++ {
		go l.mediaWorker()
	}
}

func (l *collectionLoader) loadCollection() error {
	return nil
}

//Fetches the off-chain data from the on-chain account & passes to `fileDataWorker`
func (l *collectionLoader) metaDataWorker() {

}

//Downloads required files & passes to `mediaWorker`
func (l *collectionLoader) fileDataWorker() {

}

//Stores media data down to SQL
func (l *collectionLoader) mediaWorker() {
	for m := range l.mediaIn {
		//TODO SAVE TO DB
		log.Printf("M: %s", m.MediaUri)
	}
}
