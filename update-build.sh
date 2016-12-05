#### for wfdemo-git-oem uncomment below part
sudo pkill wfdemo-git-oem

git reset --hard

rm wfdemo-git-oem
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"

go build -o wfdemo-git-oem
nohup ./wfdemo-git-oem &
