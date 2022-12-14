package main

import (
  "bump/server/server"
  "os"
  "fmt"
  "net"
  "time"
  "database/sql"
  "github.com/joho/godotenv"
  _ "github.com/go-sql-driver/mysql"
)

type RawEntry *server.RawEntry
type ValidEntry *server.ValidEntry

func push(db *sql.DB, c chan RawEntry) {
  insert, err := db.Prepare("insert into raw_point (timestamp, latitude, longitude, error, value) values (?, ?, ?, ?, ?)")
  if err != nil {
    panic(err.Error())
  }
  defer insert.Close()

  lastUpdateYear, lastUpdateMonth, lastUpdateDay := time.Now().UTC().Date()
  for entry := range c {
    currentTime := time.Now().UTC()
    currentYear, currentMonth, currentDay := currentTime.Date()
    if currentYear > lastUpdateYear || currentMonth > lastUpdateMonth || currentDay > lastUpdateDay {
      lastUpdateYear, lastUpdateMonth, lastUpdateDay = currentYear, currentMonth, currentDay
      update(db)
      currentTime = time.Now().UTC()
    }

    if (!entry.Validate(currentTime)) {
      continue
    }

    _, err = insert.Exec(time.UnixMilli(entry.Timestamp).UTC(), entry.Lat, entry.Long, entry.Err, entry.Value)
    if err != nil {
      panic(err.Error())
    }
  }
}

func update(db *sql.DB) {
  rows, err := db.Query("select timestamp, latitude, longitude, error, value from raw_point")
  if err != nil {
    panic(err.Error())
  }
  defer rows.Close()

  insert, err := db.Prepare("insert into valid_point (timestamp, place_id, value) values (?, ?, ?)")
  if err != nil {
    panic(err.Error())
  }
  defer insert.Close()

  for hasNext := true; hasNext; {
    raw := make([]RawEntry, 100)
    var i uint64
    for i, hasNext = 0, rows.Next(); i < 100 && hasNext; hasNext = rows.Next() {
      if entry := ReadRawEntry(rows); entry != nil {
        append(raw[i], entry)
        i++
      }
    }

    if err := rows.Err(); err != nil {
      panic(err.Error())
    }

    if v, err := server.ValidateRawEntries(raw); err != nil {
      panic(err.Error())
    }

    for valid := range v {
      _, err = insert.Exec(valid.Timestamp, valid.PlaceID, valid.Value)
      if err != nil {
        panic(err.Error())
      }
    }
  }

  _, err = db.Exec("delete from raw_point")
  if err != nil {
    panic(err.Error())
  }

  average(db, placeIds)
}

func average(db *sql.DB, places map[string]*string) {
  rows, err := db.Query(
    `select distinct v.place_id from valid_point v
    left join average_point a on v.place_id = a.place_id
    where a.place_id is null
    group by v.place_id`
  )
  if err != nil {
    panic(err.Error())
  }
  defer rows.Close()

  insert, err := db.Prepare("insert into average_point (place_id, route, value) values (?, ?, 0)")
  if err != nil {
    panic(err.Error())
  }
  defer insert.Close()

  for rows.Next() {
    var id string
    if err := rows.Scan(&id); err != nil {
      continue
    }
    if name, err := server.GetRoadName(id); err != nil {
      continue
    }
    insert.Exec(id, name)
  }

  if err := rows.Err(); err != nil {
    panic(err.Error())
  }
  /*
  select v3.place_id, AVG(v3.value) average from valid_point v3
  join (
    select v1.place_id, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2)) / COUNT(v1.value)) as std
    from valid_point v1
    join (
      select place_id, AVG(value) average
      from valid_point group by place_id
    ) v2 on v1.place_id = v2.place_id
    group by v1.place_id
  ) v4 on v3.place_id = v4.place_id
  where v3.value between v4.average - v4.std and v4.average + v4.std
  group by v3.place_id
  */
  _, err := db.Exec(
    `update average_point a, (
       select v3.place_id, AVG(v3.value) average from valid_point v3
       join (
         select v1.place_id, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2)) / COUNT(v1.value)) as std
         from valid_point v1
         join (
           select place_id, AVG(value) average
           from valid_point group by place_id
         ) v2 on v1.place_id = v2.place_id
         group by v1.place_id
       ) v4 on v3.place_id = v4.place_id
       where v3.value between v4.average - v4.std and v4.average + v4.std
       group by v3.place_id
     ) v set a.value = v.average where a.place_id = v.place_id`
  )
  if err != nil {
    panic(err.Error())
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
    server.ReadPacket(buffer, n, c)
  }
}
