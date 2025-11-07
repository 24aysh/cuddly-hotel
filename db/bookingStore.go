package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBooking(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}
type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("bookings"),
	}
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.D{
		{
			Key: "$set", Value: update,
		},
	}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	if err != nil {
		return err
	}
	return nil
}

func (b *MongoBookingStore) GetBooking(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	resp, err := b.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err = resp.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (b *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	err = b.coll.FindOne(ctx, bson.M{
		"_id": oid,
	}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (b *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := b.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(bson.ObjectID)
	return booking, nil

}
