#!/bin/bash
#kill `cat run.pid`
npm run build
go build
./stock-portfolio
#nohup ./stock-portfolio > logs/app.log 2>&1&
#echo $! > run.pid
