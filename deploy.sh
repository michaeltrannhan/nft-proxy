GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ./runtime/main.go
rsync -aP ./main root@159.65.58.224:/var/www/degencdn
rsync -aP ./main root@159.65.58.224:/var/www/degencdn-devnet

rsync -aP ./docs root@159.65.58.224:/var/www/degencdn/

rsync -aP ./raw_cache/solana/ root@159.65.58.224:/mnt/volume_lon1_01/cache/solana/


ssh root@159.65.58.224


sudo mount -o defaults,nofail,discard,noatime /dev/disk/by-id/scsi-0DO_Volume_volume-lon1-01 /mnt/volume_lon1_01 ext4
