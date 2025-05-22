package database

import "go.mongodb.org/mongo-driver/v2/mongo"

type CommonOperations struct {
	Client *mongo.Client
}

type Files struct {
	CreatedAt    int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt    int64  `json:"deletedAt" bson:"deletedAt"`
	IsDeleted    bool   `json:"isDeleted" bson:"isDeleted"`
	OriginalName string `json:"originalName" bson:"originalName"`
	Size         int64  `json:"size" bson:"size"`
	UID          string `json:"uid" bson:"uid"`
	UpdatedAt    int64  `json:"updatedAt" bson:"updatedAt"`
}

type Metrics struct {
	CreatedAt      int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt      int64  `json:"deletedAt" bson:"deletedAt"`
	Downloads      int64  `json:"downloads" bson:"downloads"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
	LastDownloaded int64  `json:"lastDownloaded" bson:"lastDownloaded"`
	LastViewed     int64  `json:"lastViewed" bson:"lastViewed"`
	UID            string `json:"uid" bson:"uid"`
	UpdatedAt      int64  `json:"updatedAt" bson:"updatedAt"`
	Views          int64  `json:"views" bson:"views"`
}

type Users struct {
	CreatedAt      int64  `json:"createdAt" bson:"createdAt"`
	DeletedAt      int64  `json:"deletedAt" bson:"deletedAt"`
	Email          string `json:"email" bson:"email"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
	PasswordHash   string `json:"-" bson:"passwordHash"`
	Role           string `json:"role" bson:"role"`
	SetUpCompleted bool   `json:"setUpCompleted" bson:"setUpCompleted"`
	UID            string `json:"uid" bson:"uid"`
	UpdatedAt      int64  `json:"updatedAt" bson:"updatedAt"`
}
