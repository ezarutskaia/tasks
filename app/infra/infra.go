package infra

import (
	"tasks/app/infra/repository"
)

type Infra struct {
	Repository *repository.Repository
}