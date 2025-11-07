package main

import (
	"context"
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client     *mongo.Client
	ctx        = context.Background()
	roomStore  db.RoomStore
	userStore  db.UserStoreInterface
	hotelStore db.HotelStore
)

func seedUser(isAdmin bool, email, fname, lname string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "123123123",
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(*user))

}

func seedHotel(name, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []bson.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		return err
	}
	rooms := []types.Room{
		{
			Type:  types.SinglePersonRoomType,
			Size:  "small",
			Price: 713,
		},
		{
			Type:  types.DoubleRoomType,
			Size:  "medium ",
			Price: 1031,
		},
		{
			Type:  types.DeluxeRoomType,
			Size:  "king",
			Price: 1999,
		},
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := seedHotel("Bellucia", "France", 4); err != nil {
		log.Fatal(err)
	}
	seedHotel("Angela white", "Use me", 5)
	seedHotel("Cumatoz", "Russia", 5)
	seedUser(false, "james@gmail.com", "james", "in the ass")
	seedUser(true, "angel@gmail.com", "angel", "on my dick")

}

func init() {
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
