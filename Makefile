runDev:
	 go run main.go -envString dev

runProd:
	go run main.go -envString prod

build:
	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o gin-admin-api main.go

buildMac:
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build -o gin-admin-api main.go

chmod:
	chmod 777 gin-admin-api

startDev:
	pm2 start gin-admin-api -o ./out.log -e ./error.log --log-date-format="YYYY-MM-DD HH:mm Z"

startProd:
	pm2 start gin-admin-api -o ./out.log -e ./error.log --log-date-format="YYYY-MM-DD HH:mm Z" -- -envString prod