CREATE TABLE countries(
  id integer primary key auto_increment,
  name varchar(255)
);

CREATE TABLE regions(
  id integer primary key auto_increment,
  name varchar(255),
  country_id int
);

CREATE INDEX regions_country_id ON regions (country_id);

CREATE TABLE prefectures(
  id integer primary key auto_increment,
  name varchar(255),
  region_id int
);

CREATE INDEX prefectures_region_id ON prefectures (region_id);

INSERT INTO countries (name) values ('france');
INSERT INTO regions (name, country_id) values ('Tokyo', 2);
