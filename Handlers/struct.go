package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Data struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	LocationDates []Locations
	LocationMap   []LocationPoints
}

type Locations struct {
	Country string
	City    string
	Dates   []string
}

type RelationsData struct {
	Index []struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

func geocodeAddress(LocationDates []Locations) ([]LocationPoints, error) {
	MapLocation := []LocationPoints{}
	for i := 0; i < len(LocationDates); i++ {
		apiURL := fmt.Sprintf("https://dev.virtualearth.net/REST/v1/Locations?query=%s&key=ApZ8vBLSKATZ_ByiXrOzPJ4DKmeGcsIYCL2XgmLQt8-MIrbP5fEYW7z4mTdvojwZ", LocationDates[i].City+LocationDates[i].Country)

		response, err := http.Get(apiURL)
		if err != nil {
			return nil, fmt.Errorf("geocoding request failed: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("geocoding request failed with status code: %d", response.StatusCode)
		}

		var data struct {
			ResourceSets []struct {
				Resources []struct {
					Point struct {
						Coordinates []float64 `json:"coordinates"`
					} `json:"point"`
				} `json:"resources"`
			} `json:"resourceSets"`
		}

		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse geocoding response: %v", err)
		}

		if len(data.ResourceSets) == 0 || len(data.ResourceSets[0].Resources) == 0 {
			return nil, fmt.Errorf("no results found for the address")
		}

		latitude := data.ResourceSets[0].Resources[0].Point.Coordinates[0]
		longitude := data.ResourceSets[0].Resources[0].Point.Coordinates[1]
		Location := LocationPoints{Latitude: latitude, Longitude: longitude}
		MapLocation = append(MapLocation, Location)
	}

	return MapLocation, nil
}

func FetchArtistData() ([]Data, error) {
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	relationURL := "https://groupietrackers.herokuapp.com/api/relation"

	var artists []Artist
	var relation RelationsData

	err := fetchData(artistsURL, &artists)
	if err != nil {
		return nil, err
	}

	err = fetchData(relationURL, &relation)
	if err != nil {
		return nil, err
	}

	data := []Data{}
	// Update artists with dates and locations
	for i, rel := range relation.Index {
		if len(rel.DatesLocations) > 0 {
			LocationData := []Locations{}
			for location, dates := range rel.DatesLocations {
				L := strings.Split(location, "-")
				city := L[0]
				country := L[1]
				Location := Locations{Country: country, City: city, Dates: dates}
				LocationData = append(LocationData, Location)
			}
			if i < len(artists) {
				ArtistData := Data{
					ID:            artists[i].ID,
					Image:         artists[i].Image,
					Name:          artists[i].Name,
					Members:       artists[i].Members,
					CreationDate:  artists[i].CreationDate,
					FirstAlbum:    artists[i].FirstAlbum,
					LocationDates: LocationData,
				}
				if i == 20 {
					ArtistData = Data{
						ID:            artists[i].ID,
						Image:         "https://images.virgula.me/2023/04/mamonas-assassinas.webp",
						Name:          artists[i].Name,
						Members:       artists[i].Members,
						CreationDate:  artists[i].CreationDate,
						FirstAlbum:    artists[i].FirstAlbum,
						LocationDates: LocationData,
					}

				}
				data = append(data, ArtistData)
			}
		}
	}

	return data, nil
}
