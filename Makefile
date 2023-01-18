runDev:
	 ENV=dev go run main.go
dockerProdUp:
	docker-compose -f docker-compose-prod.yml up