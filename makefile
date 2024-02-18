run:
	docker-compose up --build cloud-concierge

fumpt:
	gofumpt -l -w .
