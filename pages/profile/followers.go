package profile

import (
	"net/http"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

// GetFollowers shows the followers of a particular user.
func GetFollowers(ctx *aero.Context) string {
	nick := ctx.Get("nick")
	viewUser, err := arn.GetUserByNick(nick)

	if err != nil {
		return ctx.Error(http.StatusNotFound, "User not found", err)
	}

	followers := viewUser.Followers()
	arn.SortUsersLastSeen(followers)

	return ctx.HTML(components.ProfileFollowers(followers, viewUser, utils.GetUser(ctx), ctx.URI()))

}
