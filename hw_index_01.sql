CREATE TABLE flights (
    flight_id INT PRIMARY KEY AUTO_INCREMENT,
    airline VARCHAR(100),
    flight_number VARCHAR(20),
    departure_airport VARCHAR(100),
    arrival_airport VARCHAR(100),
    departure_time DATETIME,
    arrival_time DATETIME,
    price DECIMAL(10, 2)
);

CREATE TABLE customers (
    customer_id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100),
    phone_number VARCHAR(20)
);

CREATE TABLE bookings (
    booking_id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT,
    flight_id INT,
    booking_date DATETIME,
    seat VARCHAR(50),
    class INT,
);


CREATE VIEW FlightsWithBookingCountToday AS
SELECT 
    f.flight_id,
    f.airline,
    f.flight_number,
    f.departure_airport,
    f.arrival_airport,
    f.departure_time,
    f.arrival_time,
    f.price,
    COUNT(b.booking_id) AS booking_count_today,
    IFNULL(COUNT(b.booking_id) * f.price, 0) AS total_price_today
FROM 
    flights f
LEFT JOIN 
    bookings b ON f.flight_id = b.flight_id
WHERE 
    DATE(b.booking_date) = CURDATE()
GROUP BY 
    f.flight_id;


-- create index
CREATE INDEX idx_booking_customer_flight ON bookings (customer_id, flight_id, booking_date);
CREATE INDEX idx_booking_class ON bookings (class);
CREATE INDEX idx_flights ON flights (flight_number);
CREATE INDEX idx_flights_departure_arrival_airport ON flights (departure_airport, arrival_airport);
CREATE INDEX idx_flights_departure_arrival_time ON flights (departure_time, arrival_time);