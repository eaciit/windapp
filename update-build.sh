#### for wfdemo-git-oem uncomment below part

git checkout .

sudo pkill wfdemo-oem
rm wfdemo-oem
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"

go build -o wfdemo-oem
nohup ./wfdemo-oem &