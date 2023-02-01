package main

import "time"

type TwitchClips struct {
	Data       Data               `json:"data"`
	Extensions ResponseExtensions `json:"extensions"`
}
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	Typename    string `json:"__typename"`
}
type Curator struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"displayName"`
	Typename    string `json:"__typename"`
}
type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"boxArtURL"`
	Typename  string `json:"__typename"`
}
type Broadcaster struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageURL"`
	PrimaryColorHex string `json:"primaryColorHex"`
	Typename        string `json:"__typename"`
}
type Clip struct {
	ID              string      `json:"id"`
	Slug            string      `json:"slug"`
	URL             string      `json:"url"`
	EmbedURL        string      `json:"embedURL"`
	Title           string      `json:"title"`
	ViewCount       int         `json:"viewCount"`
	Language        string      `json:"language"`
	Curator         Curator     `json:"curator"`
	Game            Game        `json:"game"`
	Broadcaster     Broadcaster `json:"broadcaster"`
	ThumbnailURL    string      `json:"thumbnailURL"`
	CreatedAt       time.Time   `json:"createdAt"`
	DurationSeconds int         `json:"durationSeconds"`
	ChampBadge      interface{} `json:"champBadge"`
	Typename        string      `json:"__typename"`
}
type ClipNode struct {
	Cursor   interface{} `json:"cursor"`
	Node     Clip        `json:"node"`
	Typename string      `json:"__typename"`
}
type Clips struct {
	PageInfo PageInfo   `json:"pageInfo"`
	Edges    []ClipNode `json:"edges"`
	Typename string     `json:"__typename"`
}
type User struct {
	ID       string `json:"id"`
	Clips    Clips  `json:"clips"`
	Typename string `json:"__typename"`
}
type Data struct {
	User User `json:"user"`
}
type ResponseExtensions struct {
	DurationMilliseconds int    `json:"durationMilliseconds"`
	OperationName        string `json:"operationName"`
	RequestID            string `json:"requestID"`
}
