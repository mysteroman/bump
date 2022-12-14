package server

import "math"

type ValidEntry struct {
  PlaceID string
  Timestamp int64
  Value float64
}

func (entry RawEntry) Validate(placeId string, latitude, longitude float64) *ValidEntry {
  if distanceBetween(latitude, longitude, entry.Lat, entry.Long) > entry.Err {
    return nil
  }
  return &ValidEntry{placeId, entry.Timestamp, entry.Value}
}

func distanceBetween(lat1, long1, lat2, long2 float64) float64 {
  const radius float64 = 6378000

  lat1 = toRad(lat1)
  lat2 = toRad(lat2)
  long1 = toRad(long1)
  long2 = toRad(long2)

  dlat := math.Sin((lat2 - lat1) / 2)
  dlat *= dlat

  dlong := math.Sin((long2 - long1) / 2)
  dlong *= dlong

  result := dlat + math.Cos(lat1) * math.Cos(lat2) * dlong
  result = math.Sqrt(result)
  return 2 * radius * math.Asin(result)
}

func toRad(v float64) float64 {
  return math.Pi * v / 180
}