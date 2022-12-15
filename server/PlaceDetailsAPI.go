package main

import (
  "net/http"
  "io"
  "encoding/json"
)

const base_url = "https://maps.googleapis.com/maps/api/place/details/json?fields=address_components&key=AIzaSyCpzHxqHktEbM3YTgTzRHZi6ilSJZdtoKc&place_id="

type response struct {
  result struct {
    address_components []struct {
      types []string
      long_name string
    }
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

  return find(data.result.address_components, "route") + ", " + find(data.result.address_components, "locality"), nil
}

func find(components []struct{types []string; long_name string}, category string) string {
  for _, component := range components {
    for _, cat := range component.types {
      if cat == category {
        return component.long_name
      }
    }
  }
  return category
}