package app

import (
	"bufio"
	"os"
	"vyshkat/cli"
	"vyshkat/repository"
	"vyshkat/service"
)

type App struct {
	handler *cli.CLIHandler
}

func NewApp() *App {
	container := BuildContainer()
	return &App{
		handler: Resolve[*cli.CLIHandler](container),
	}
}

func (a *App) Run() {
	a.handler.HandleRequests()
}

func BuildContainer() *Container {
	container := NewContainer()

	Register[service.InventoryRepository](container, func() service.InventoryRepository {
		return repository.NewMemoryInventoryRepository()
	})

	Register[service.ServiceCenter](container, func() service.ServiceCenter {
		return service.ItemInspector{}
	})

	Register[*service.InventoryService](container, func() *service.InventoryService {
		repo := Resolve[service.InventoryRepository](container)
		inspector := Resolve[service.ServiceCenter](container)
		return service.NewInventoryService(repo, inspector)
	})

	Register[*cli.CLIHandler](container, func() *cli.CLIHandler {
		scanner := bufio.NewScanner(os.Stdin)
		service := Resolve[*service.InventoryService](container)
		return cli.NewCLIHandler(service, scanner, os.Stdout)
	})
	return container
}
