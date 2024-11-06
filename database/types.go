package database

type File struct {
	CreatedAt    int64  `json:"createdAt"`
	Downloads    int64  `json:"downloads"`
	OriginalName string `json:"originalName"`
	UID          string `json:"uid"`
	Size         int    `json:"size"`
}
