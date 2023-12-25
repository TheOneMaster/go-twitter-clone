run:
	@./tailwind -i ./templates/layouts/main.css -o ./static/main.css
	@go run ./main.go
