ENV_FILE=.env
MIGRATIONS_FOLDER=migrations
PROTO_FOLDER=protobuf/proto
BUF_VERSION?=1.58.0

env-example:
	awk -F'=' 'BEGIN {OFS="="} \
    	/^[[:space:]]*#/ {print; next} \
    	/^[[:space:]]*$$/ {print ""; next} \
    	NF>=1 {gsub(/^[[:space:]]+|[[:space:]]+$$/, "", $$1); print $$1"="}' .env > .env.example
	echo ".env.example generated successfully."

mongo-mock-up:
	go run ./migrations/mongo/migrate.go up
