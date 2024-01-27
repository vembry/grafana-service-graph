default:
	echo "hi there!"

start:	
	make down
	make up
	for i in `seq 1 100`; do curl --location 'http://localhost:3000/ping' --header 'Content-Type: application/json' --data '{"hosts": ["http://app-go-1","http://app-go-2","http://app-go"]}'; done

up:
	docker compose -f "compose.yml" up -d --build --remove-orphans

down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)

watch:
	docker compose alpha watch