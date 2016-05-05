CREATE TABLE countries(
  "id" serial primary key,
  "name" varchar(255)
);

CREATE TABLE regions(
  "id" serial primary key,
  "name" varchar(255),
  "country_id" int,
  "order" int
);

CREATE INDEX regions_country_id ON regions ("country_id");

CREATE TABLE prefectures(
  "id" serial primary key,
  "name" varchar(255),
  "region_id" int
);

CREATE INDEX prefectures_region_id ON prefectures ("region_id");

INSERT INTO countries ("name") values ('france');
INSERT INTO regions ("name", "country_id", "order") values ('Tokyo', 2, 10);
