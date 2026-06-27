package service

import (
	"zhanbu/config"
	"zhanbu/internal/model"
)

// ApplyDefaultPromptProfile attaches default profile metadata to a record when
// the divination type is explicitly bound to an enabled prompt profile.
func ApplyDefaultPromptProfile(record *model.DivinationRecord, profiles *config.ProfilesConfig) {
	if record == nil || profiles == nil {
		return
	}

	profileID, ok := profiles.DefaultBindings[record.Type]
	if !ok {
		return
	}

	profile, err := profiles.DefaultProfile(record.Type)
	if err != nil {
		return
	}

	record.PromptProfileID = profileID
	record.PromptProfileName = profile.Name
	record.PromptProfileVersion = profile.Version
}
