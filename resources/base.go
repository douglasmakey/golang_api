package resources

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type Error struct {
	Data    interface{} `json:"error"`
}
