package main

import (
	"bufio"
	"os"
	"vyshkat/cli"
	"vyshkat/repository"
	"vyshkat/service"
)

func main() {
	repo := repository.NewMemoryInventoryRepository()
	service := service.NewInventoryService(repo)

	scanner := bufio.NewScanner(os.Stdin)
	handler := cli.NewCLIHandler(service, scanner, os.Stdout)

	handler.HandleRequests()
}
