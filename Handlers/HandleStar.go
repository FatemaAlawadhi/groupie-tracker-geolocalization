package Handlers

import (
	"encoding/json"
	"groupie-tracker-filter/Error"
	"html/template"
	"net/http"
	"strconv"
)

func HandleStarsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ArtistData,err := FetchArtistData()
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error sug")
			return
		}

		suggestions, err := SearchSuggestions(ArtistData)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error sug")
			return
		}

		suggestionsJSON, err := json.Marshal(suggestions)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error4")
			return
		}

		//Parse and Excute
		tmpl, err := template.ParseFiles("WebPages/Stars/Stars.html", "WebPages/Stars/Card.html", "WebPages/Stars/Filters.html")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error3")
			return
		}
	
		data := struct {
			Artists       []Data
			SuggestionsJS template.JS
		}{
			Artists:       ArtistData,
			SuggestionsJS: template.JS(suggestionsJSON),
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error Ex")
			return
		}
	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed 5")
		return
	}
}

func SearchSuggestions(data []Data) ([]string, error) {
	var Suggestions []string
	var txt string
	for i, artist := range data {
		
		txt = artist.Name + " - artist/band"
		Suggestions = append(Suggestions, txt)

		for _, member := range artist.Members {
			txt = member + " - member"
			Suggestions = append(Suggestions, txt)
		}

		txt = strconv.Itoa(artist.CreationDate) + " - Creation Date"
		Suggestions = append(Suggestions, txt)

		txt = artist.FirstAlbum + " - First Album"
		Suggestions = append(Suggestions, txt)
		
		for _, dl := range data[i].LocationDates{
			txt = dl.City + "-" + dl.Country + " - Location"
			Suggestions = append(Suggestions, txt)
			for _, date := range dl.Dates{
				txt = date + " - Date"
				Suggestions = append(Suggestions, txt)
			}
		}
	}
	return Suggestions, nil
}


