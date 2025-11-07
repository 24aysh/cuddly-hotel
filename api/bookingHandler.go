package api

import (
	"fmt"
	"hotel-reservation/db"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBooking(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return fmt.Errorf("cannot cancel the booking")
	}
	update := bson.M{
		"cancelled": true,
	}
	if err = h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), update); err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"Message": "Booking cancelled",
	})
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return fmt.Errorf("errer here")
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"Error": "not authorized",
		})
	}

	return c.JSON(booking)
}
