package main

import (
	"errors"

	"github.com/cloudfoundry-community/cfenv"
)

//GetUserProvidedServiceByName retreives a user provided service by name
func GetUserProvidedServiceByName(name string, env *cfenv.App) (*cfenv.Service, error) {
	ups := env.Services["user-provided"]
	for _, s := range ups {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, errors.New("Cannot find service named " + name)
}
