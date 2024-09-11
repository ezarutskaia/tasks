package app

import (
	"tasks/app/domain"
	"tasks/app/infra"
	"tasks/app/interfaces"

)

type App struct {
	Domain *domain.Domain
	Infra *infra.Infra
	Interfaces *interfaces.Interfaces
}