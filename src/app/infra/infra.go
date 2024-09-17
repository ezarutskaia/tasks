package infra

import (
	"tasks/src/app/infra/repository"
)

type Infra struct {
	Repository *repository.Repository
}