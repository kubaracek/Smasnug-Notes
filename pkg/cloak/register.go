package cloak

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type RegisterValues map[string]map[string]string

// Updates registers and returns previous values as new RegisterValues
func (r RegisterValues) Update() (RegisterValues, error) {
	prevRegisters := make(RegisterValues)

	// Get all registers first to store previous values
	for path, registers := range r {
		// Initialize the inner map for each path in prevRegisters
		if _, exists := prevRegisters[path]; !exists {
			prevRegisters[path] = make(map[string]string)
		}

		for name, _ := range registers {
			value, err := getRegister(path, name)
			if err != nil {
				return nil, err
			}

			prevRegisters[path][name] = value
		}
	}

	// Set new registers
	for path, registers := range r {
		for name, value := range registers {
			err := setRegister(path, name, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return prevRegisters, nil
}

type RegisterCloak struct {
	updateRegisters RegisterValues
}

func NewRegisterCloak(registers RegisterValues) *RegisterCloak {
	return &RegisterCloak{updateRegisters: registers}
}

func (r RegisterCloak) CloakExecution(fn func() error) error {
	prevRegisters, err := r.updateRegisters.Update()

	// Revert registers back
	defer prevRegisters.Update()

	err = fn()
	if err != nil {
		return fmt.Errorf("while being cloaked: %v", err)
	}

	return nil
}

func getRegister(path, name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	defer key.Close()
	if err != nil {
		return "", err
	}

	value, _, err := key.GetStringValue(name)
	return value, err
}

func setRegister(path, name, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
	defer key.Close()
	if err != nil {
		return err
	}

	return key.SetStringValue(name, value)
}
