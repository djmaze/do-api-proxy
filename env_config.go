package main

import (
  "log"
  "net/http"
  "os"
  "regexp"
  "strings"
)

type EnvConfig struct {
  ListenPort string
  ApiToken string
  TargetApiURL string
  TargetApiToken string
  Volumes []DOVolume
}

type DOVolume struct {
  Id string
  Name string
}

func GetEnvConfig() EnvConfig {
  config := EnvConfig{
    ListenPort: getEnvWithFallback("LISTEN_PORT", "1338"),
    ApiToken: getEnv("API_TOKEN"),
    TargetApiURL: getEnvWithFallback("TARGET_API_URL", "https://api.digitalocean.com"),
    TargetApiToken: getEnv("TARGET_API_TOKEN"),
  }

  ids := strings.Split(getEnv("VOLUMES"), ",")

  log.Printf("Getting volumes from DO..")
  client := InitClient(config.TargetApiToken)
  volumes := client.getVolumes(ids)
  config.Volumes = volumes

  return config
}

func (config EnvConfig) IsRequestAuthCorrect(req *http.Request) bool {
  authorization := req.Header.Get("Authorization")
  re := regexp.MustCompile(`Bearer (.+)`)
  matches := re.FindStringSubmatch(authorization)
  token := matches[1]

  if token == "" {
    log.Fatalf("Request does not contain token\n")
    return false
  }

  if config.ApiToken != token {
    log.Printf("Given token is not correct")
    return false
  }

  return true
}

func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
    log.Fatalf("ENV variable for %s missing!", key)
    return ""
  }
}
