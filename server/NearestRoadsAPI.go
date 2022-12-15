package main

import (
  "strings"
  "net/http"
  "io"
  "encoding/json"
)

const base_roads_url = "https://roads.googleapis.com/v1/nearestRoads?key=AIzaSyCpzHxqHktEbM3YTgTzRHZi6ilSJZdtoKc&points="

type point struct {
  location struct { latitude, longitude float64 }
  originalIndex uint64
  placeId string
}

type roads_response struct {
  snappedPoints []point
}

func ValidateRawEntries(entries []*RawEntry) ([]*ValidEntry, error) {
  if len(entries) > 100 {
    panic("Too many raw entries")
  }

  url := buildUrl(entries)
  res, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  var data roads_response
  if err := json.Unmarshal(body, &data); err != nil {
    return nil, err
  }

  result := make([]*ValidEntry, len(data.snappedPoints))
  var i uint64 = 0
  for _, point := range data.snappedPoints {
    if valid_point := entries[point.originalIndex].Validate(point.placeId, point.location.latitude, point.location.longitude); valid_point != nil {
      result[i] = valid_point
      i++
    }
  }

  return result, nil
}

func buildUrl(entries []*RawEntry) string {
  url := base_roads_url
  points := make([]string, 0)
  for _, entry := range entries {
    if entry == nil {
      break
    }
    points = append(points, entry.LocationString())
  }
  return url + strings.Join(points, "%7C")
}