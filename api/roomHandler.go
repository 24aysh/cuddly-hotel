package api

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoomHandler struct {
	store *db.Store
}
type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPerson"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}

func (r *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := r.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (r *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return nil
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID := c.Params("id")
	roomOid, err := bson.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"Error": "Internal Server error",
		})
	}

	avail, err := r.isRoomAvailable(c.Context(), roomOid, params)
	if err != nil {
		return err
	}
	if !avail {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"Error": "Room already booked within this period",
		})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOid,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	fmt.Printf("%+v\n", booking)
	inserted, err := r.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (r *RoomHandler) isRoomAvailable(ctx context.Context, roomID bson.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := r.store.Booking.GetBooking(ctx, where)
	if err != nil {
		return false, nil
	}
	if len(bookings) > 0 {
		return false, nil
	}
	return true, nil
}
