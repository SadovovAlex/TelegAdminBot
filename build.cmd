SET CGO_ENABLED=1
echo Build the Go application
cd ./cmd/bot
go build -v -o ../../distr/telegbotadmin.exe 
cd ../../
copy .env distr\.env

cd ./distr
start.cmd

pause
