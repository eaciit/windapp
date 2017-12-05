sudo pkill wfdemo-git-oem

git reset --hard

rm wfdemo-git-oem
sudo mv nohup.out /mnt/data/ostrorealtime/applog/log_oem_$(date +"%Y%m%d_%H%M%S")

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"

go build -o wfdemo-git-oem
sudo nohup ./wfdemo-git-oem &
