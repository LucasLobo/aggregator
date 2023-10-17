PROJECT_NAME := aggregator-cli

build: $(PROJECT_NAME)
$(PROJECT_NAME):
	go build -o $(PROJECT_NAME) .

run-default:
	go run . moving-average --window_size 10 --input_file data/events.json --output_folder data/output

example: build
	./$(PROJECT_NAME) moving-average --window_size 10 --input_file data/events.json --output_folder data/output

clean-output:
	rm $(CURDIR)/data/output/events_*.json || true