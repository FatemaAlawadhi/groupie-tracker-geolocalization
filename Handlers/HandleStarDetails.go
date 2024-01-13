package Handlers

import (
	"groupie-tracker-filter/Error"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type LocationPoints struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func HandleStarDetailsPage(w http.ResponseWriter, r *http.Request) {
	// Fetch the main data
	idstr := strings.TrimPrefix(r.URL.Path, "/stardetails/")
	ArtistsData, err := FetchArtistData()
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if idint <= 0 || idint >= 53 {
		Error.RenderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}

	ArtistData := ArtistsData[idint-1]
	MapLocation := []LocationPoints{}
	MapLocation, err = geocodeAddress(ArtistData.LocationDates)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ArtistData = Data{
		ID:            ArtistData.ID,
		Image:         ArtistData.Image,
		Name:          ArtistData.Name,
		Members:       ArtistData.Members,
		CreationDate:  ArtistData.CreationDate,
		FirstAlbum:    ArtistData.FirstAlbum,
		LocationDates: ArtistData.LocationDates,
		LocationMap:   MapLocation,
	}

	tmpl, err := template.ParseFiles("WebPages/StarDetails.html")
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = tmpl.Execute(w, ArtistData)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
