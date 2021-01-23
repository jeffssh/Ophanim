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
	Pipe        string `yaml:"pipe"`
	Command     string `yaml:"command"`
	Enabled     bool   `yaml:"enabled"`
	cmd         *exec.Cmd
}

// Start - Start the module per the specified command
func (m *Module) Start() (err error) {
	if !m.Enabled {
		return
	}
	/*
		c := winio.PipeConfig{
			//D:PAI(A;;0x100116;;;WD) - Write access to Everyone
			SecurityDescriptor: "D:PAI(A;;0x100116;;;WD)",
			MessageMode:        true,
		}
		l, err := winio.ListenPipe(`\\.\pipe\`+m.Pipe, &c)
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
	*/
	if m.Command == "" {
		log.Printf("Module %s started\n", m.Name)
		return
	}
	cmd := exec.Command(m.Command)
	cmd.Env = append(os.Environ())
	err = cmd.Start()
	if err == nil {
		m.cmd = cmd
		go func() {
			cmd.Wait()
			log.Printf("Module %s command %s exited with code: %d\n", m.Name, m.Command, cmd.ProcessState.ExitCode())
			m.cmd = nil
		}()
	} else {
		return
	}
	log.Printf("Module %s started with command %s\n", m.Name, m.Command)
	return
}

// Stop - Stop the module cmd
func (m *Module) Stop() (err error) {
	if !m.Enabled {
		return
	}
	if m.cmd == nil {
		log.Printf("Module %s stopped\n", m.Name)
		return
	}

	err = m.cmd.Process.Kill()
	if err == nil {
		log.Printf("Module %s command %s stopped\n", m.Name, m.Command)
	}
	return
}

func (m *Module) String() string {
	return fmt.Sprintf("Name: %s\nDescription: %s\nPipe: %s\nCommand: %s\nEnabled: %v", m.Name, m.Description, m.Pipe, m.Command, m.Enabled)
}

// LoadModule - load a module from a yaml file specified by modulePath
func LoadModule(modulePath string) (m *Module, err error) {
	m = &Module{}
	data, err := ioutil.ReadFile(modulePath)
	if err != nil {
		log.Printf("error reading module yaml file: %v\n", err)
		return
	}

	err = yaml.Unmarshal(data, m)
	if err != nil {
		log.Printf("error unmarshalling yaml: %v\n", err)
		return
	}
	return
}

// LoadAllModules - load all modules in a directory
func LoadAllModules(moduleDir string) (modules []*Module) {
	moduleFiles, err := ioutil.ReadDir(moduleDir)
	if err != nil {
		log.Printf("Error loading module files: %v\n", err)
		return
	}

	for _, f := range moduleFiles {
		m, err := LoadModule(moduleDir + f.Name())
		if err == nil {
			modules = append(modules, m)
		}
	}
	return
}
