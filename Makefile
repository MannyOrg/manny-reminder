all: provide_command

I=""

provide_command:
	@echo "Please specify command"

generate_mocks_all:
	docker run -v "${PWD}":/src -w /src vektra/mockery --all

generate_mocks:
	docker run -v "${PWD}":/src -w /src vektra/mockery --name=${I} --recursive=true
