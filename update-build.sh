#### for wfdemo-git-prod uncomment below part

sudo pkill wfdemo-git-dev

git reset --hard

rm wfdemo-git-dev
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-prod"

go build -o wfdemo-git-dev
nohup ./wfdemo-git-dev &
