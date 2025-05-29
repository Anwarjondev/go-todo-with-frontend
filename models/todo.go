package models

// Todo represents a single todo item
type Todo struct {
	Id        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"user_id"`
}

// View represents the data structure for the template
type View struct {
	Todos []Todo `json:"todos"`
}
