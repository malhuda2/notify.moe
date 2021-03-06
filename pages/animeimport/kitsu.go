package animeimport

import (
	"fmt"
	"net/http"

	"github.com/animenotifier/kitsu"
	"github.com/fatih/color"

	"github.com/animenotifier/arn"

	"github.com/aerogo/aero"
	"github.com/animenotifier/notify.moe/utils"
)

// Kitsu anime import.
func Kitsu(ctx *aero.Context) string {
	id := ctx.Get("id")
	user := utils.GetUser(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
	}

	kitsuAnimeObj, err := arn.Kitsu.Get("Anime", id)

	if kitsuAnimeObj == nil {
		return ctx.Error(http.StatusNotFound, "Kitsu anime not found", err)
	}

	kitsuAnime := kitsuAnimeObj.(*kitsu.Anime)

	// Convert
	anime, characters, relations, episodes := arn.NewAnimeFromKitsuAnime(kitsuAnime)

	// Add user ID to the anime
	anime.CreatedBy = user.ID

	// Save in database
	anime.Save()
	characters.Save()
	relations.Save()
	episodes.Save()

	// Log
	fmt.Println(color.GreenString("✔"), anime.ID, anime.Title.Canonical)

	return ""
}
