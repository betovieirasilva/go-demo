package entity

//vars
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
	Stock  int64   `json:"stock" binding:"stockValidator"`
}
