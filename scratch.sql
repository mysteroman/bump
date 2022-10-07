use bump;

create table raw_point (
    id bigint not null primary key auto_increment,
    timestamp datetime not null,
    latitude double not null,
    longitude double not null,
    error double not null,
    value double not null
);

create table valid_point (
    id bigint not null primary key auto_increment,
    timestamp datetime not null,
    place_id varchar(255) not null,
    value double not null
);

create table average_point (
      place_id varchar(255) not null primary key,
      route varchar(255) not null,
      value double not null
);
