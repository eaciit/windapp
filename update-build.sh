#### for wfdemo-git-oem uncomment below part

git reset --hard

sudo pkill wfdemo-git-oem
rm wfdemo-git-oem
rm nohup.out

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"

go build -o wfdemo-git-oem
nohup ./wfdemo-git-oem &
