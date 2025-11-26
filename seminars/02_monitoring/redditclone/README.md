To build and run the app inside container:

docker compose exec -it redditclone /bin/bash
cd /app
go build -o redditclone ./...
./redditclone > app-logs/redditclone.log 2>&1
