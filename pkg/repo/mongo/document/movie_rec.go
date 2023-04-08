package document

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Movie struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	MovieID   int                `bson:"movie_id"`
	OverView  string             `bson:"over_view,omitempty"`
	Url       string             `bson:"url,omitempty"`
	ImageUrl  string             `bson:"image_url,omitempty"`
	LeadActor string             `bson:"lead_actor"`
	Tags      []string           `bson:"tags"`
	CreateTs  time.Time          `bson:"create_ts"`
	UpdateTs  time.Time          `bson:"update_ts"`
}
