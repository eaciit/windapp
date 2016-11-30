#### for wfdemo-git-prod uncomment below part

git checkout .

sudo pkill wfdemo-prod
rm wfdemo-prod
rm nohup.out

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-prod"

go build -o wfdemo-prod
nohup ./wfdemo-prod &