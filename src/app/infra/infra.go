package infra

import (
	"tasks/src/app/infra/pdf"
	"tasks/src/app/infra/repository"
)

type Infra struct {
	Repository *repository.Repository
	Pdf *pdf.Pdf
}