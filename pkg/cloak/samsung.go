package cloak

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

const REGISTER_PATH = `HARDWARE\DESCRIPTION\System\BIOS`
const REGISTER_PRODUCT_NAME = "SystemProductName"
const REGISTER_PRODUCT_MANUFACTURER = "SystemManufacturer"
const SAMSUNG_LAPTOP_PRODUCT_NAME = "NP960XFG-KC4UK"
const SAMSUNG_LAPTOP_MANUFACTURER_NAME = "Samsung"

type Samsung struct{}

func NewSamsung() *Samsung {
	return &Samsung{}
}

func (r Samsung) CloakExecuting(fn func() error) error {
	// Store previous SystemProductName and SystemManufacturer
	productName, err := r.getRegister(REGISTER_PRODUCT_NAME)
	if err != nil {
		return fmt.Errorf("failed to read SystemProductName: %v", err)
	}

	manufacturer, err := r.getRegister(REGISTER_PRODUCT_MANUFACTURER)
	if err != nil {
		return fmt.Errorf("failed to read SystemManufacturer: %v", err)
	}

	// revert values back
	defer func(r Samsung, productName, productManufacturer string) {
		err := r.revert(productName, productManufacturer)
		if err != nil {
			panic("failed to revert registry values back to their original values")
		}
	}(r, productName, manufacturer)

	fmt.Println("Cloaking PC as a Samsung Device")
	err = r.setRegister(REGISTER_PRODUCT_NAME, SAMSUNG_LAPTOP_PRODUCT_NAME)
	if err != nil {
		return fmt.Errorf("failed to write to registry to change product name: %v", err)
	}

	err = r.setRegister(REGISTER_PRODUCT_NAME, SAMSUNG_LAPTOP_MANUFACTURER_NAME)
	if err != nil {
		return fmt.Errorf("failed to write to registry to change manufacturer name: %v", err)
	}

	err = fn()
	if err != nil {
		return fmt.Errorf("while being cloaked: %v", err)
	}

	return nil
}

func (r Samsung) revert(productName, productManufacturer string) error {
	fmt.Println("Reverting registry to original values")
	var err error
	err = r.setRegister("SystemProductName", productName)
	if err != nil {
		return fmt.Errorf("failed to restore SystemProductName: %v", err)
	}
	err = r.setRegister("SystemManufacturer", productManufacturer)
	if err != nil {
		return fmt.Errorf("failed to restore SystemManufacturer: %v", err)
	}
	return nil
}

// Read a cloak key value
func (r Samsung) getRegister(name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, REGISTER_PATH, registry.QUERY_VALUE)
	defer key.Close()
	if err != nil {
		return "", err
	}

	value, _, err := key.GetStringValue(name)
	return value, err
}

// Set a cloak key value
func (r Samsung) setRegister(name, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, REGISTER_PATH, registry.SET_VALUE)
	defer key.Close()
	if err != nil {
		return err
	}

	return key.SetStringValue(name, value)
}
