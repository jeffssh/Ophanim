package module

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// Module - Struct defined by module yaml files, contains all information necessary to launch and communicate with the module
type Module struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Key         string `yaml:"key"`
	Command     string `yaml:"command"`
	Enabled     bool   `yaml:"enabled"`
	cmd         *exec.Cmd
}

// Map - Map of Modules
type Map map[string]*Module

// Start - Start the module per the specified command
func (m *Module) Start() (err error) {
	if !m.Enabled {
		return
	}
	if m.Command == "" {
		log.Printf("Module %s started", m.Name)
		return
	}
	cmd := exec.Command(m.Command)
	cmd.Env = append(os.Environ())
	err = cmd.Start()
	if err == nil {
		m.cmd = cmd
		go func() {
			cmd.Wait()
			log.Printf("Module %s command %s exited with code: %d", m.Name, m.Command, cmd.ProcessState.ExitCode())
			m.cmd = nil
		}()
	} else {
		return
	}
	log.Printf("Module %s started with command %s", m.Name, m.Command)
	return
}

// Stop - Stop the module cmd
func (m *Module) Stop() (err error) {
	if !m.Enabled {
		return
	}
	if m.cmd == nil {
		log.Printf("Module %s stopped", m.Name)
		return
	}

	err = m.cmd.Process.Kill()
	if err == nil {
		log.Printf("Module %s command %s stopped", m.Name, m.Command)
	}
	return
}

func (m *Module) String() string {
	return fmt.Sprintf("Name: %s\nDescription: %s\nKey: %s\nCommand: %s\nEnabled: %v", m.Name, m.Description, m.Key, m.Command, m.Enabled)
}

// LoadModule - load a module from a yaml file specified by modulePath
func LoadModule(modulePath string) (m *Module, err error) {
	m = &Module{}
	data, err := ioutil.ReadFile(modulePath)
	if err != nil {
		log.Printf("error reading module yaml file: %v", err)
		return
	}

	err = yaml.Unmarshal(data, m)
	if err != nil {
		log.Printf("error unmarshalling yaml: %v", err)
		return
	}
	return
}

// LoadAllModules - load all modules in a directory
func LoadAllModules(moduleDir string) (modules Map) {
	moduleFiles, err := ioutil.ReadDir(moduleDir)
	if err != nil {
		log.Printf("Error loading module files: %v", err)
		return
	}

	modules = make(Map)
	for _, f := range moduleFiles {
		m, err := LoadModule(moduleDir + f.Name())
		if err == nil {
			modules[m.Key] = m
		}
	}
	return
}

// Stop - stop all modules in a map
func (modules Map) Stop() {
	for _, m := range modules {
		err := m.Stop()
		if err != nil {
			fmt.Printf("Error when stopping module %s command %s: %v\n", m.Name, m.Command, err)
		}
	}
}

// Start - start all modules in a map
func (modules Map) Start() {
	for _, m := range modules {
		err := m.Start()
		if err != nil {
			fmt.Printf("Error when starting module %s, command %s: %v\n", m.Name, m.Command, err)
		}
	}
}

// Enable - enable all modules in a map
func (modules Map) Enable() {
	for _, m := range modules {
		m.Enabled = true
	}
	fmt.Printf("Enabled %d modules\n", len(modules))
}

// Disable - enable all modules in a map
func (modules Map) Disable() {
	for _, m := range modules {
		m.Enabled = false
	}
	fmt.Printf("Disabled %d modules\n", len(modules))
}

// Describe - start all modules in a map
func (modules Map) Describe() {
	for _, m := range modules {
		sep := "======================================="
		fmt.Printf("%s\n%v\n%s\n", sep, m, sep)
	}
}
