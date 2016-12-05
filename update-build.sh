#### for wfdemo-git-prod uncomment below part

git checkout .

sudo pkill wfdemo-git-prod
rm wfdemo-git-prod
rm nohup.out

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-prod"

go build -o wfdemo-git-prod
nohup ./wfdemo-git-prod &
