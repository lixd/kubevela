package system

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloud-native-application/rudrx/api/types"
)

const defaultVelaHome = ".vela"
const VelaHomeEnv = "VELA_HOME"

func GetVelaHomeDir() (string, error) {
	if custom := os.Getenv(VelaHomeEnv); custom != "" {
		return custom, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultVelaHome), nil
}

func GetCapCenterDir() (string, error) {
	home, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "centers"), nil
}

func GetRepoConfig() (string, error) {
	home, err := GetCapCenterDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "config.yaml"), nil
}

func GetCapabilityDir() (string, error) {
	home, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "capabilities"), nil
}

func GetEnvDir() (string, error) {
	homedir, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, "envs"), nil
}

func GetCurrentEnvPath() (string, error) {
	homedir, err := GetVelaHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, "curenv"), nil
}

func InitDirs() error {
	if err := InitCapabilityDir(); err != nil {
		return err
	}
	if err := InitCapCenterDir(); err != nil {
		return err
	}
	return nil
}

func InitCapCenterDir() error {
	home, err := GetCapCenterDir()
	if err != nil {
		return err
	}
	return StatAndCreate(filepath.Join(home, ".tmp"))
}

func InitCapabilityDir() error {
	dir, err := GetCapabilityDir()
	if err != nil {
		return err
	}
	return StatAndCreate(dir)
}

func GetApplicationDir(envName string) (string, error) {
	appDir := filepath.Join(GetEnvDirByName(envName), "applications")
	return appDir, StatAndCreate(appDir)
}

const EnvConfigName = "config.json"

func InitDefaultEnv() error {
	envDir, err := GetEnvDir()
	if err != nil {
		return err
	}
	defaultEnvDir := filepath.Join(envDir, types.DefaultEnvName)
	StatAndCreate(defaultEnvDir)
	data, _ := json.Marshal(&types.EnvMeta{Namespace: types.DefaultAppNamespace, Name: types.DefaultEnvName})
	if err = ioutil.WriteFile(filepath.Join(defaultEnvDir, EnvConfigName), data, 0644); err != nil {
		return err
	}
	curEnvPath, err := GetCurrentEnvPath()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(curEnvPath, []byte(types.DefaultEnvName), 0644); err != nil {
		return err
	}
	return nil
}

func StatAndCreate(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func GetEnvDirByName(name string) string {
	envdir, _ := GetEnvDir()
	return filepath.Join(envdir, name)
}
