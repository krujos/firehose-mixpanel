package main

import (
	"errors"
	"log"

	"github.com/cloudfoundry-community/cfenv"
)

//GetUserProvidedServiceByName retreives a user provided service by name
func GetUserProvidedServiceByName(name string, env *cfenv.App) (*cfenv.Service, error) {
	log.Println(len(env.Services))

	ups := env.Services["user-provided"]
	for _, s := range ups {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, errors.New("Cannot find service named " + name)
}
