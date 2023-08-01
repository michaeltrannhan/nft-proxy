GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ./runtime/main.go
rsync -aP ./main root@159.65.58.224:/var/www/degencdn
rsync -aP ./main root@159.65.58.224:/var/www/degencdn-devnet

rsync -aP ./docs root@159.65.58.224:/var/www/degencdn/


ssh root@159.65.58.224