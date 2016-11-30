#### for wfdemo-git-dev uncomment below part

git checkout .

sudo pkill wfdemo-git-dev
rm wfdemo-git-dev
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-dev"

go build -o wfdemo-git-dev
nohup ./wfdemo-git-dev &
