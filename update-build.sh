#### for wfdemo-git-dev uncomment below part

sudo pkill wfdemo-dev
rm wfdemo-dev
rm nohup.out

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git"

git pull

$GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git"


go build -o wfdemo-dev
nohup ./wfdemo-dev &




# #### for wfdemo-git-oem uncomment below part

# sudo pkill wfdemo-dev
# rm wfdemo-dev
# rm nohup.out

# $GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git-oem" -to="eaciit/wfdemo-git"

# git pull

# $GOPATH/bin/gorep -path="." -from="eaciit/wfdemo-git" -to="eaciit/wfdemo-git-oem"


# go build -o wfdemo-dev
# nohup ./wfdemo-dev &
