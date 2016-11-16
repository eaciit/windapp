sudo pkill wfdemo-receiver

$GOPATH/bin/gorep -path="." -from="github.com/eaciit/windapp-dev" -to="github.com/eaciit/windapp"

svn update ../../../library/
svn update ../../watcher

$GOPATH/bin/gorep -path="." -from="github.com/eaciit/windapp" -to="github.com/eaciit/windapp-dev"

rm wfdemo-receiver
rm nohup.out

go build -o wfdemo-receiver
nohup ./wfdemo-receiver &
