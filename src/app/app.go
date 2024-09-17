package app

import (
	"tasks/src/app/domain"
	"tasks/src/app/infra"
	"tasks/src/app/interfaces"

)

type App struct {
	Domain *domain.Domain
	Infra *infra.Infra
	Interfaces *interfaces.Interfaces
}