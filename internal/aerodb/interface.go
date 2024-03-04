package aerodb

import "time"

type Sqlite3DB interface {
    OpenDB(fname string) (error)
    CloseDB() (error) 
    PlanTrip(trip Trip) (TripID int, err error)
    EndTrip(tripID int) (error)
    GetTrips(from, to string) ([]Trip, error)
    GetAllTrips() ([]Trip, error)
    TakeSeat(tripID int, passenger string, seat int) (error)
    GetFreeSeats(tripID int) ([]int, error)
    AddCompany(name string) (error)
    DelCompany(name string) (error)
    AddPlane(name, companyName string, seats int) (error)
    DelPlane(name string) (error)
    AddPassenger(name string) (error)
}


type Trip struct {
    id int
    company, plane int
    timeOut, timeIn time.Time
    townOut, townIn int
}