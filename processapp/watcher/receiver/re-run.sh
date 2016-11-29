sudo pkill wfdemo-git-watcher

rm wfdemo-git-watcher
rm nohup.out

go build -o wfdemo-git-watcher
nohup ./wfdemo-git-watcher &
