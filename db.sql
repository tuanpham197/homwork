CREATE TABLE air_ports (
  id BIGINT UNSIGNED auto_increment primary key NOT NULL,
  code varchar(100) NOT NULL,
  name varchar(100) NOT NULL
)
    ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci;

INSERT INTO airports (code, name)
VALUES ('AIR_001', 'Airport 01'),
       ('AIR_002', 'Airport 02'),
       ('AIR_003', 'Airport 03'),
       ('AIR_004', 'Airport 04'),
       ('AIR_005', 'Airport 05'),
       ('AIR_006', 'Airport 06');
