package database

type File struct {
	CreatedAt    int64  `json:"createdAt"`
	Downloads    int64  `json:"downloads"`
	OriginalName string `json:"originalName"`
	Size         int64  `json:"size"`
	UID          string `json:"uid"`
}
