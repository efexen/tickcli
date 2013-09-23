package main

import (
  "os"
  "os/user"
  "io/ioutil"
  "encoding/json"
  "path/filepath"
)

const DotFileName = ".tickrc"

type Config struct {
  Company string
  Email string
  Password string
}

func (config *Config) WriteToFile() error {
  bytes, err := json.Marshal(config)
  if err != nil { return err }

  file, err := os.Create(config.FilePath())
  if err != nil { return err }

  file.Write(bytes[:])
  file.Close()

  return nil
}

func (config *Config) ReadFromFile() {
  bytes, _ := ioutil.ReadFile(config.FilePath())
  json.Unmarshal(bytes, &config)
}

func (config *Config) FilePath() (path string) {
  user, _ := user.Current()
  return filepath.Join(user.HomeDir, DotFileName)
}
