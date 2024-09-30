package pdf

type TaskDTO struct {
	Id     int		`json:"id"`
	Title  string	`json:"title"`
	User   string	`json:"user"`
	Weight int		`json:"weight"`
}