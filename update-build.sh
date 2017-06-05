sudo pkill wfdemo-git-oem

git reset --hard

rm wfdemo-git-oem
sudo rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"

go build -o wfdemo-git-oem
sudo nohup ./wfdemo-git-oem &
