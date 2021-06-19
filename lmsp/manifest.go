package lmsp

import "time"

type Manifest struct {
	Type          string    `json:"type"` // always "word-blocks" for EV3?
	AutoDelete    bool      `json:"autoDelete"`
	Created       time.Time `json:"created"`
	ID            string    `json:"id"`
	LastSaved     time.Time `json:"lastsaved"`
	Size          int       `json:"size"` // always 0?
	Name          string    `json:"name"`
	SlotIndex     int       `json:"slotIndex"`
	WorkspaceX    float64   `json:"workspaceX"`
	WorkspaceY    float64   `json:"workspaceY"`
	ZoomLevel     float64   `json:"zoomLevel"`
	ShowAllBlocks bool      `json:"showAllBlocks"`
	Version       int       `json:"version"` // always 5 for EV3?
	Hardware      map[string]struct {
		Name               string `json:"name"` // name of EV3 brick
		Connection         string `json:"connection"`
		LastConnectedHubID string `json:"lastConnectedHubId"`
		ID                 string `json:"id"`
		Type               string `json:"type"`
	} `json:"hardware"`
	Extensions []string `json:"extensions"` // ev3events, ev3move, ev3motor, ev3sensors
	State      struct {
		PlayMode        string `json:"playMode"`
		CanvasDrawerTab string `json:"canvasDrawerTab"`
	} `json:"state"`
}
