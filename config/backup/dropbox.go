package backup

import (
	"github.com/BurntSushi/toml"
)

// Dropbox is dropboxc backup folder setting
type Dropbox struct {
	Src       string
	Dst       string
	BackupNum int
}

// NewDropbox return Dropbox from yaml data
func NewDropbox(settingText string) *Dropbox {
	var d Dropbox

	_, err := toml.Decode(settingText, &d)
	if err != nil {
		return nil
	}
	return &d
}
