sudo pkill wfdemo-receiver

$GOPATH/bin/gorep -path="." -from="eaciit/ostrowfm-dev" -to="eaciit/ostrowfm"

svn update ../../../library/
svn update ../../watcher

$GOPATH/bin/gorep -path="." -from="eaciit/ostrowfm" -to="eaciit/ostrowfm-dev"

rm wfdemo-receiver
rm nohup.out

go build -o wfdemo-receiver
nohup ./wfdemo-receiver &
