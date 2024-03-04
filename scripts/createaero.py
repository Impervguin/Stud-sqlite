import os
import sqlite3
import database as data

FILENAME = "./database/aero.sqlite3"

try:
    os.mkdir("database")
except FileExistsError:
    pass

if (os.path.exists(FILENAME)):
    os.remove(FILENAME)

db = sqlite3.connect("database/aero.sqlite3")

db.execute("""CREATE TABLE IF NOT EXISTS Trip (
           id INTEGER primary key,
           company_id INTEGER,
           plane_id INTEGER,
           time_out TIMESTAMP,
           time_in TIMESTAMP,
           town_out VARCHAR,
           town_in VARCHAR
)""")

db.execute("""CREATE TABLE IF NOT EXISTS Plane (
           id INTEGER primary key,
           name VARCHAR UNIQUE,
           company_id INTEGER,
           seats INTEGER
)""")

db.execute("""CREATE TABLE IF NOT EXISTS Company (
           id INTEGER primary key,
           name VARCHAR UNIQUE
)""")

db.execute("""CREATE TABLE IF NOT EXISTS Passenger (
           id INTEGER primary key,
           name VARCHAR UNIQUE
)""")

db.execute("""CREATE TABLE IF NOT EXISTS Taken (
           id INTEGER primary key,
           trip_id INTEGER,
           passenger_id INTEGER,
           place INTEGER
)""")




db.executemany("INSERT INTO Passenger(name) VALUES (?)", data.PASSENGERS)
db.executemany("INSERT INTO Company(name) VALUES (?)", data.COMPANIES)
db.executemany("""INSERT INTO Plane(name, company_id, seats) 
               SELECT ?, id, ?
               From Company
               where name=?
               """, data.PLANES)
db.executemany("""INSERT INTO Trip(company_id, plane_id, time_out, time_in, town_out, town_in)
               VALUES ((SELECT id from Company where name=?),
               (SELECT id from Plane where name=?),
               ?, ?, ?, ?)
""", data.TRIPS)
db.executemany("""INSERT INTO Taken(trip_id, passenger_id, place)
               SELECT ?, id, ?
               From Passenger
               where name=?
               """, data.TAKEN)
db.commit()
db.close()