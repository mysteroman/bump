package main

import (
  "fmt"
  "math"
  "time"
)

// Raw entry received from UDP
type RawEntry struct {
  Timestamp int64
  Lat, Long, Err, Value float64
}

func ReadPacket(data []byte, n int64, c chan *RawEntry) {
  if n % 40 != 8 {
    return
  }
  n = n / 40
  key := readInt64(data[0:8])
  timestamp := readInt64(data[8:16])
  if !verifyKey(timestamp, key) {
    return
  }
  var i int64
  for i = 0; i < n; i++ {
    c <- newRawEntry(data[i*40:(i+1)*40])
  }
}

func newRawEntry(data []byte) *RawEntry {
  return &RawEntry{
    readInt64(data[0:8]),
    readFloat64(data[8:16]),
    readFloat64(data[16:24]),
    readFloat64(data[24:32]),
    readFloat64(data[32:40])}
}

func ReadRawEntry(reader interface{Scan(dest ...any) error}) *RawEntry {
  var entry RawEntry
  if err := reader.Scan(&entry.Timestamp, &entry.Lat, &entry.Long, &entry.Err, &entry.Value); err != nil {
    return nil
  }
  return &entry
}

func (entry RawEntry) IsValid(currentTime time.Time) bool {
  if currentMillis := currentTime.UnixMilli(); entry.Timestamp > currentMillis || entry.Timestamp + 1000 < currentMillis {
    return false
  }
  if math.IsNaN(entry.Lat) || math.IsInf(entry.Lat, 0) || entry.Lat < -90 || entry.Lat > 90 {
    return false
  }
  if math.IsNaN(entry.Long) || math.IsInf(entry.Long, 0) || entry.Long < -180 || entry.Long > 180 {
    return false
  }
  if math.IsNaN(entry.Err) || math.IsInf(entry.Err, 0) || entry.Err < 0 || entry.Err > 50 {
    return false
  }
  if math.IsNaN(entry.Value) || math.IsInf(entry.Value, 0) || entry.Value < 0 {
    return false
  }
  return true
}

func (entry RawEntry) LocationString() string {
  return fmt.Sprintf("%f,%f", entry.Lat, entry.Long)
}

func verifyKey(timestamp, key int64) bool {
  var hash int64 = 0
  for i := 0; i < 64; i++ {
    hash *= 31
    if (timestamp & (1 << i)) != 0 {
      hash++
    }
  }
  return hash == key
}

func readUInt64(bytes []byte) uint64 {
  var result uint64 = 0
  for i, b := range bytes {
    if i >= 8 {
      break
    }
    result = (result << 8) | uint64(b)
  }
  return result
}

func readInt64(bytes []byte) int64 {
  return int64(readUInt64(bytes))
}

func readFloat64(bytes []byte) float64 {
  return math.Float64frombits(readUInt64(bytes))
}