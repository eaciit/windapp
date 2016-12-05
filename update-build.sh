#### for wfdemo-git-prod uncomment below part

sudo pkill wfdemo-git-prod

git reset --hard

rm wfdemo-git-prod
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-prod"

go build -o wfdemo-git-prod
nohup ./wfdemo-git-prod &
