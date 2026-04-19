package main

import (
	"log"
	"net/http"

	"github.com/PaleBlueDot1990/cinema-seat-booking/internal/adapters/redis"
	"github.com/PaleBlueDot1990/cinema-seat-booking/internal/booking"
	"github.com/PaleBlueDot1990/cinema-seat-booking/internal/utils"
)

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

var movies = []movieResponse {
	{ID: "dune", Title: "Dune Part Three", Rows: 5, SeatsPerRow: 8},
	{ID: "doom", Title: "Avengers Doomsday", Rows: 4, SeatsPerRow: 6},
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("GET /movies", listMovies)

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeats)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}
