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
	inspector := service.ItemInspector{}
	service := service.NewInventoryService(repo, inspector)

	scanner := bufio.NewScanner(os.Stdin)
	handler := cli.NewCLIHandler(service, scanner, os.Stdout)

	handler.HandleRequests()
}
