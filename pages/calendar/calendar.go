package calendar

import (
	"sort"
	"time"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/arn/validator"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

var weekdayNames = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

// Get ...
func Get(ctx *aero.Context) string {
	user := utils.GetUser(ctx)
	oneWeek := 7 * 24 * time.Hour

	now := time.Now()
	year, month, day := now.Date()
	now = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// Weekday index that we start with, Sunday is 0.
	weekdayIndex := int(now.Weekday())

	// Create days
	days := make([]*utils.CalendarDay, 7, 7)

	for i := 0; i < 7; i++ {
		days[i] = &utils.CalendarDay{
			Name:    weekdayNames[(weekdayIndex+i)%7],
			Entries: []*utils.CalendarEntry{},
		}
	}

	// Add anime episodes to the days
	for animeEpisodes := range arn.StreamAnimeEpisodes() {
		for _, episode := range animeEpisodes.Items {
			if !validator.IsValidDate(episode.AiringDate.Start) {
				continue
			}

			// Since we validated the date earlier, we can ignore the error value.
			airingDate, _ := time.Parse(time.RFC3339, episode.AiringDate.Start)

			// Subtract from the starting date offset.
			since := airingDate.Sub(now)

			// Ignore entries in the past and more than 1 week away.
			if since < 0 || since >= oneWeek {
				continue
			}

			dayIndex := int(since / (24 * time.Hour))

			entry := &utils.CalendarEntry{
				Anime:   animeEpisodes.Anime(),
				Episode: episode,
				Class:   "calendar-entry mountable",
			}

			if user != nil {
				animeListItem := user.AnimeList().Find(entry.Anime.ID)

				if animeListItem != nil && (animeListItem.Status == arn.AnimeListStatusWatching || animeListItem.Status == arn.AnimeListStatusPlanned) {
					entry.Class += " calendar-entry-personal"
				}
			}

			days[dayIndex].Entries = append(days[dayIndex].Entries, entry)
		}
	}

	for i := 0; i < 7; i++ {
		sort.Slice(days[i].Entries, func(a, b int) bool {
			airingA := days[i].Entries[a].Episode.AiringDate.Start
			airingB := days[i].Entries[b].Episode.AiringDate.Start

			if airingA == airingB {
				return days[i].Entries[a].Anime.Title.Canonical < days[i].Entries[b].Anime.Title.Canonical
			}

			return airingA < airingB
		})
	}

	return ctx.HTML(components.Calendar(days, user))
}