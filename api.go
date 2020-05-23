package main

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/godo"
)

type Client struct {
  Token string
  GodoClient *godo.Client
}

func InitClient(token string) (Client) {
  godoClient := godo.NewFromToken(token)
  return Client{GodoClient: godoClient}
}

func (client *Client) getVolumes(names []string) []DOVolume {
  ctx := context.TODO()

  apiVolumes, _, err := client.GodoClient.Storage.ListVolumes(
    ctx,
    &godo.ListVolumeParams{
      ListOptions: &godo.ListOptions{Page: 1, PerPage: 100}, // FIXME Allow more than 100 volumes
    },
  )
  if err != nil {
    log.Fatalf("%v", err)
  }

  doVolumes := make([]DOVolume, len(names))
  for _, name := range names {
    found := false
    for _, apiVolume := range apiVolumes {
      if strings.Compare(apiVolume.Name, name) == 0 {
        doVolumes = append(doVolumes, DOVolume{
          Id: apiVolume.ID,
          Name: apiVolume.Name,
        })
        found = true
        break
      }
    }
    if !found {
      log.Fatalf("Could not find volume %s\n", name)
    }
  }

  return doVolumes
}
