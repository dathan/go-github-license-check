package sheets

import (
	"google.golang.org/api/sheets/v4"
)

type Service struct {
	srv          *sheets.Service
	currentSheet string
}
