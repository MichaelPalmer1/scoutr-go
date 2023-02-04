package models

// History : History object
type History struct {
	Time string `json:"time" firestore:"time"`
	Data Record `json:"data" firestore:"data"`
}
