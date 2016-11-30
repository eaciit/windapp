#### for wfdemo-git-prod uncomment below part

git checkout .

sudo pkill wfdemo-git-dev
rm wfdemo-git-dev
rm nohup.out

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-prod"

go build -o wfdemo-git-dev
nohup ./wfdemo-git-dev
