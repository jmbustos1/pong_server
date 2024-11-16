.PHONY: go postgres dockerup api

go:
	sudo docker exec -it pong_client_container /bin/bash

postgres:
	docker exec -it conergie-postgres /bin/bash

api:
	docker exec -it conergie-api-1 /bin/bash

dockerup:
	sudo docker compose up
