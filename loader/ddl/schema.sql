CREATE TABLE stations
(
    id        varchar not null primary key,
    name      varchar not null,
    longitude double precision,
    latitude  double precision
);

CREATE TABLE line_has_stations
(
    line_name  varchar not null,
    station_id varchar not null,
    primary key (line_name, station_id)
);
CREATE INDEX IDX_LINE_HAS_STATIONS_STATION_ID ON line_has_stations USING btree (station_id);

select * from stations;
select * from lines;
select * from line_has_stations;

--- Query lines by station
SELECT lhs.line_name
FROM stations
         LEFT JOIN line_has_stations lhs on stations.id = lhs.station_id
WHERE name = 'Oxford Circus';

--- Query station name by line name
SELECT s.name
FROM line_has_stations
         LEFT JOIN stations s on line_has_stations.station_id = s.id
WHERE line_name = 'Victoria';