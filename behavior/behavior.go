package behavior

import (
	"github.com/imkira/go-observer/v2"
)

type Behavior struct {
	Preferences observer.Property[Preferences]
}

func NewBehavior() (*Behavior, error) {
	prefs, err := LoadPreferences()
	if err != nil {
		return nil, err
	}
	prefs.Defaults()

	return &Behavior{
		Preferences: observer.NewProperty(*prefs),
	}, nil
}
