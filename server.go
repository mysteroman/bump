package main

import (
  "os"
  "fmt"
  "net"
  "math"
  "time"
  "database/sql"
  "github.com/joho/godotenv"
  _ "github.com/go-sql-driver/mysql"
)

type RawEntry struct {
  Timestamp int64
  Lat, Long, Err, Value float64
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

func readFloat64(bytes []byte) float64 {
  return math.Float64frombits(readUInt64(bytes))
}

func newRawEntry(data []byte) RawEntry {
  return RawEntry{
    int64(readUInt64(data[0:8])),
    readFloat64(data[8:16]),
    readFloat64(data[16:24]),
    readFloat64(data[24:32]),
    readFloat64(data[32:40])}
}

func push(db *sql.DB, c chan RawEntry) {
  insert, err := db.Prepare("insert into raw_point (timestamp, latitude, longitude, error, value) values (?, ?, ?, ?, ?)")
  if err != nil {
    panic(err.Error())
  }
  defer insert.Close()

  for entry := range c {
    currentTime := time.Now().UTC()


    currentMillis := currentTime.UnixMilli()
    if entry.Timestamp > currentMillis || entry.Timestamp + 1000 < currentMillis {
      continue
    }
    if math.IsNaN(entry.Lat) || math.IsInf(entry.Lat, 0) || entry.Lat < -90 || entry.Lat > 90 {
      continue
    }
    if math.IsNaN(entry.Long) || math.IsInf(entry.Long, 0) || entry.Long < -180 || entry.Long > 180 {
      continue
    }
    if math.IsNaN(entry.Err) || math.IsInf(entry.Err, 0) || entry.Err < 0 || entry.Err > 1000 {
      continue
    }
    if math.IsNaN(entry.Value) || math.IsInf(entry.Value, 0) || entry.Value < 0 {
      continue
    }

    _, err := insert.Exec(time.UnixMilli(entry.Timestamp).UTC(), entry.Lat, entry.Long, entry.Err, entry.Value)
    if err != nil {
      panic(err.Error())
    }
  }
}

func main() {
  godotenv.Load()

  PORT := ":53"
  s, err := net.ResolveUDPAddr("udp4", PORT)
  if err != nil {
    panic(err.Error())
  }

  conn, err := net.ListenUDP("udp4", s)
  if err != nil {
    panic(err.Error())
  }

  defer conn.Close()

  dsn := os.Getenv("DSN")
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()

  db.SetConnMaxLifetime(time.Minute * 3)
  db.SetMaxOpenConns(10)
  db.SetMaxIdleConns(10)

  err = db.Ping()
  if err != nil {
    panic(err.Error())
  }

  fmt.Println("Awaiting requests...")

  buffer := make([]byte, 1000)
  c := make(chan RawEntry, 1)
  defer close(c)

  go push(db, c)

  for {
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil {
      panic(err.Error())
    }
    if n % 40 != 0 {
      continue
    }
    n = n / 40
    for i := 0; i < n; i++ {
      c <- newRawEntry(buffer[i*40:(i+1)*40])
    }
  }
}
