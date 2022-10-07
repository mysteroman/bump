create table raw_points (
    id int not null primary key auto_increment,
    timestamp timestamp not null,
    latitude double not null,
    longitude double not null,
    error double not null,
    value double not null
);

create table valid_points (
    id int not null primary key auto_increment,
    timestamp timestamp not null,
    place_id varchar(255) not null,
    value double not null
);

create table average_points (
      place_id varchar(255) not null primary key,
      route varchar(255) not null,
      value double not null
);
