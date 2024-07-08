package main

type GameEvent struct {
	Version          string `json:"version,omitempty"`
	LeaderboardQueue string `json:"leaderboard_queue"`
	Serial           string `json:"serial,omitempty"`
	Model            string `json:"model,omitempty"`
	AutomaticUpdates bool   `json:"automatic_updates,omitempty"`
}
