package database

type Files struct {
	CreatedAt    int64  `json:"createdAt" bson:"createdAt"`
	OriginalName string `json:"originalName" bson:"originalName"`
	Size         int64  `json:"size" bson:"size"`
	UID          string `json:"uid" bson:"uid"`
}

type Metrics struct {
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	Downloads int64  `json:"downloads" bson:"downloads"`
	UID       string `json:"uid" bson:"uid"`
	Views     int64  `json:"views" bson:"views"`
}
