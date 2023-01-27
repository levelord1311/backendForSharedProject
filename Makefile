run_api:
	cd ./api_service/app/; go run ./cmd/main/app.go

run_user:
	cd ./user_service/app/; go run ./cmd/main/app.go

run_lot:
	cd ./lot_service/app/; go run ./cmd/main/app.go
