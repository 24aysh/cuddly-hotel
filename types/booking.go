package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Booking struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     bson.ObjectID `bson:"userID" json:"userID"`
	RoomID     bson.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons int           `bson:"numPerson" json:"numPerson"`
	FromDate   time.Time     `bson:"fromDate" json:"fromData"`
	TillDate   time.Time     `bson:"tillDate" json:"tillDate"`
	Cancelled  bool          `bson:"cancelled" json:"cancelled"`
}


