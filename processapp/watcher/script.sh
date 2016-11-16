sudo pkill wfdemo-receiver

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-dev" -to="eaciit/wfdemo"

svn update ../../../library/
svn update ../../watcher

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo" -to="eaciit/wfdemo-dev"

rm wfdemo-receiver
rm nohup.out

go build -o wfdemo-receiver
nohup ./wfdemo-receiver &
