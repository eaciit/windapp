sudo pkill wfdemo-receiver

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git-dev" -to="eaciit/wfdemo-git"

svn update ../../../library/
svn update ../../watcher

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-dev"

rm wfdemo-receiver
rm nohup.out

go build -o wfdemo-receiver
nohup ./wfdemo-receiver &
