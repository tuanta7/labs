check:
	sudo ss -tulpn | grep 6060

run-goroutine-leak:
	go run ./cmd/goroutine-leak

k6-goroutine-leak:
	k6 run --env BASE_URL=http://localhost:6060 k6/goroutineLeak.js