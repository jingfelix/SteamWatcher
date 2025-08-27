package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/dstgo/steamapi"
	"github.com/jingfelix/SteamWatcher/configs"
	"github.com/lkzc19/bark.sdk.go"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Error(err)
	}

	client, err := steamapi.New(config.SteamAPIKey)
	if err != nil {
		log.Error(err)
	}

	// Parse user IDs from config
	steamIDs := config.SteamUserIDs
	if len(steamIDs) == 0 {
		log.Error("No Steam user IDs found in config")
		return
	}

	// Track last known online status
	lastOnlineStatus := make(map[string]bool)

	// Periodic checking function
	checkFriendsStatus := func() {
		// Get player summary
		playerSummary, err := client.ISteamUser().GetPlayerSummaries(
			strings.Join(steamIDs, ","),
		)
		if err != nil {
			log.Errorf("Error checking status for Steam ID %s: %v", steamIDs, err)
			return
		}

		onlineCount := 0

		for _, player := range playerSummary.Response.Players {
			isCurrentlyOnline := player.PersonaState > 0

			// Compare with last known status and log if changed
			if lastOnlineStatus[player.SteamId] != isCurrentlyOnline && isCurrentlyOnline {
				log.Infof("Steam user %s is now online", player.SteamId)

				// Print player struct

				onlineCount++
				req := bark.Req{
					DeviceKey: config.DeviceKey,
					Title:     "Steam User Online!",
					// Content:   fmt.Sprintf("Steam user %s is now online", player.PersonName),
					Icon: player.Avatar,
				}
				err := bark.Notify(req)
				if err != nil {
					log.Errorf("Error sending notification for Steam ID %s: %v", player.SteamId, err)
				}
			} else {
				log.Infof("Steam user %s is now offline", player.SteamId)
			}

			// Update last known status
			lastOnlineStatus[player.SteamId] = isCurrentlyOnline
		}

		if onlineCount == 0 {
			log.Infof("No one is online")
		}
	}

	// Run initial check
	checkFriendsStatus()

	// Set up periodic checking every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Block main goroutine and check status periodically
	for range ticker.C {
		checkFriendsStatus()
	}
}
