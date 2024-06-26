package middleware

import (
	"github.com/bwmarrin/discordgo"

	"{[.RootPackage]}/internal/apiserver/store"
)

func IsAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	user := i.User
	if i.User == nil {
		user = i.Member.User
	}

	admin, _ := store.S.BotAdmins().IsAdmin(Ctx, user.ID)
	if !admin {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Forbidden",
			},
		})
	}

	return admin
}
