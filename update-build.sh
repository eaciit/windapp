#### for wfdemo-git-dev uncomment below part

git checkout .

sudo pkill wfdemo-dev
rm wfdemo-dev
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git"


go build -o wfdemo-dev
nohup ./wfdemo-dev &