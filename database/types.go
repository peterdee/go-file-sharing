package database

type Files struct {
	CreatedAt    int64  `json:"createdAt"`
	OriginalName string `json:"originalName"`
	Size         int64  `json:"size"`
	UID          string `json:"uid"`
}

type Metrics struct {
	CreatedAt int64  `json:"createdAt"`
	Downloads int64  `json:"downloads"`
	UID       string `json:"uid"`
	Views     int64  `json:"views"`
}
