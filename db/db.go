package db

const (
	DBNAME = "hotel-reservation"
	DBURI  = "mongodb://localhost:27017"
)

type Store struct {
	User    UserStoreInterface
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
