package users

import (
	"sort"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

// Active ...
func Active(ctx *aero.Context) string {
	users := arn.FilterUsers(func(user *arn.User) bool {
		return user.HasAvatar() && user.HasNick() && user.IsActive()
	})

	followCount := arn.UserFollowerCountMap()

	sort.Slice(users, func(i, j int) bool {
		if users[i].HasAvatar() != users[j].HasAvatar() {
			if users[i].HasAvatar() {
				return true
			}

			return false
		}

		followersA := followCount[users[i].ID]
		followersB := followCount[users[j].ID]

		if followersA == followersB {
			return users[i].Nick < users[j].Nick
		}

		return followersA > followersB
	})

	return ctx.HTML(components.Users(users))
}

// ActiveNoAvatar ...
func ActiveNoAvatar(ctx *aero.Context) string {
	users := arn.FilterUsers(func(user *arn.User) bool {
		return user.IsActive() && !user.HasAvatar()
	})

	followCount := arn.UserFollowerCountMap()

	sort.Slice(users, func(i, j int) bool {
		if users[i].HasAvatar() != users[j].HasAvatar() {
			if users[i].HasAvatar() {
				return true
			}

			return false
		}

		followersA := followCount[users[i].ID]
		followersB := followCount[users[j].ID]

		if followersA == followersB {
			return users[i].Nick < users[j].Nick
		}

		return followersA > followersB
	})

	return ctx.HTML(components.Users(users))
}

// Osu ...
func Osu(ctx *aero.Context) string {
	users := arn.FilterUsers(func(user *arn.User) bool {
		return user.HasAvatar() && user.HasNick() && user.IsActive() && user.Accounts.Osu.PP > 0
	})

	// Sort by pp
	sort.Slice(users, func(i, j int) bool {
		return users[i].Accounts.Osu.PP > users[j].Accounts.Osu.PP
	})

	if len(users) > 50 {
		users = users[:50]
	}

	return ctx.HTML(components.OsuRankingList(users))
}

// Overwatch ...
func Overwatch(ctx *aero.Context) string {
	users := arn.FilterUsers(func(user *arn.User) bool {
		return user.HasAvatar() && user.HasNick() && user.IsActive() && user.Accounts.Overwatch.SkillRating > 0
	})

	// Sort by Skill Ratings
	sort.Slice(users, func(i, j int) bool {
		return users[i].Accounts.Overwatch.SkillRating > users[j].Accounts.Overwatch.SkillRating
	})

	if len(users) > 50 {
		users = users[:50]
	}

	return ctx.HTML(components.OverwatchRankingList(users))
}

// Staff ...
func Staff(ctx *aero.Context) string {
	users := arn.FilterUsers(func(user *arn.User) bool {
		return user.HasAvatar() && user.HasNick() && user.IsActive() && user.Role != ""
	})

	// Make order deterministic
	sort.Slice(users, func(i, j int) bool {
		return users[i].Nick < users[j].Nick
	})

	admins := &utils.UserList{
		Name:  "Developer",
		Users: []*arn.User{},
	}

	contributors := &utils.UserList{
		Name:  "Contributors",
		Users: []*arn.User{},
	}

	// contributors.Users = append(contributors.Users, )

	editors := &utils.UserList{
		Name:  "Editors",
		Users: []*arn.User{},
	}

	for _, user := range users {
		if user.Role == "admin" {
			admins.Users = append(admins.Users, user)
			continue
		}

		if user.ID == "VJOK1ckvx" || user.ID == "SUQOAFFkR" {
			contributors.Users = append(contributors.Users, user)
			continue
		}

		if user.Role == "editor" {
			editors.Users = append(editors.Users, user)
			continue
		}
	}

	userLists := []*utils.UserList{
		admins,
		contributors,
		editors,
	}

	return ctx.HTML(components.UserLists(userLists) + components.StaffRecruitment())
}
