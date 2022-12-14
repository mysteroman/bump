package server

import (
  "strings"
  "net/http"
  "io"
  "encoding/json"
)

const base_url = "https://maps.googleapis.com/maps/api/place/details/json?fields=name&key=AIzaSyCpzHxqHktEbM3YTgTzRHZi6ilSJZdtoKc&place_id="

type response struct {
  result struct {
    name string
  }
}

func GetRoadName(placeId string) (string, error) {
  url := base_url + placeId
  res, err := http.Get(url)
  if err != nil {
    return "", err
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if err != nil {
    return "", err
  }

  var data response
  if err := json.Unmarshal(body, &data); err != nil {
    return "", err
  }

  return data.result.name, nil
}