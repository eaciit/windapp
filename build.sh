sudo pkill wfdemo-oem

rm wfdemo-oem
rm nohup.out

go build -o wfdemo-oem
nohup ./wfdemo-oem &

