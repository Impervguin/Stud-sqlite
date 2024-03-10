package aerodb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

func CreateTrip(id, companyID, planeID int, timeOut, timeIn time.Time, townOut, townIn string) (Trip) {
	return Trip{id, companyID, planeID, timeOut, timeIn, townOut, townIn}
}

type AeroDB struct {
	db *sql.DB
	opened bool
}

func (db *AeroDB) OpenDB(fname string) (error) {
	// fname = fname + "?parseTime=true"
	tdb, err :=	sql.Open("sqlite3", fname)
	if (err != nil) {
		return ErrFile
	}
	db.opened = true
	db.db = tdb
	return nil
}

func (db *AeroDB) CloseDB() (error) {
	if (!db.opened) {
		return ErrNotOpened
	}
	db.db.Close()
	db.opened = false
	return nil
}


func (db *AeroDB) AddPassenger(name string) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}
	

	_, err := db.db.Exec("INSERT INTO Passenger(name) VALUES (?)", name)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if (sqliteErr.Code == sqlite3.ErrConstraint) {
			return ErrAlreadyIn
		}
	}
	return nil
}

func (db *AeroDB) AddCompany(name string) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	_, err := db.db.Exec("INSERT INTO Company(name) VALUES (?)", name)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if (sqliteErr.Code == sqlite3.ErrConstraint) {
			return ErrAlreadyIn
		}
	}
	return nil	
}

func (db *AeroDB) DelCompany(name, heritant string) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	row := db.db.QueryRow("SELECT id FROM Company WHERE name=?", name)
	if row.Scan() == sql.ErrNoRows {
		return ErrNotFound;
	}

	row = db.db.QueryRow("SELECT id FROM Company WHERE name=?", heritant)
	if row.Scan() == sql.ErrNoRows {
		return ErrNotFound;
	}

	_, err := db.db.Exec(`UPDATE Plane 
						  SET company_id=(SELECT id FROM Company WHERE name=?)
	 					  WHERE company_id=(SELECT id FROM Company WHERE name=?)`, heritant, name)

	if err != nil {
		return ErrDB;
	}

	_, err = db.db.Exec(`UPDATE Trip 
						  SET company_id=(SELECT id FROM Company WHERE name=?)
	 					  WHERE company_id=(SELECT id FROM Company WHERE name=?)`, heritant, name)

	if err != nil {
		return ErrDB;
	}

	_, err = db.db.Exec("DELETE FROM Company WHERE name=?", name)

	if err != nil {
		return ErrDB;
	}

	return nil
}

func (db *AeroDB) AddPlane(name, companyName string, seats int) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	if (seats <= 0) {
		return ErrSeatRange;
	}

	row := db.db.QueryRow("SELECT id FROM Company WHERE name=?", companyName)
	if row.Scan() == sql.ErrNoRows {
		return ErrNotFound;
	}

	_, err := db.db.Exec("INSERT INTO Plane(name, company_id, seats) VALUES (?, (SELECT id FROM Company WHERE name=?), ?)", name, companyName, seats)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if (sqliteErr.Code == sqlite3.ErrConstraint) {
			return ErrAlreadyIn
		}
	}
	return nil	
}

func filterNulls(arr []int) ([]int) {
	filt := []int{}
	for _, el := range arr {
		if (el != 0) {
			filt = append(filt, el)
		}
	}
	return filt
}

func (db *AeroDB) GetFreeSeats(tripID int) ([]int, error) {
	if (!db.opened) {
		return []int{} , ErrNotOpened
	}

	row := db.db.QueryRow("SELECT plane_id FROM Trip WHERE id=?", tripID)
	var pid int
	if row.Scan(&pid) == sql.ErrNoRows {
		return []int{}, ErrNotFound;
	}

	row = db.db.QueryRow("SELECT seats FROM Plane WHERE id=?", pid)
	var seats int
	if row.Scan(&seats) == sql.ErrNoRows {
		return []int{}, ErrNotFound;
	}

	res := make([]int, seats)

	for i:=1; i < seats + 1; i++ {
		res[i - 1] = i;
	}

	rows, err := db.db.Query("SELECT place FROM Taken WHERE trip_id=?", tripID)
	if (err != nil) {
		return []int{}, ErrDB
	}

	defer rows.Close()

	for rows.Next() {
		var taken int
		
		if err = rows.Scan(&taken); err != nil {
			return []int{}, ErrDB
		}
		fmt.Println(taken)
		res[taken  - 1] = 0
	}

	return filterNulls(res), nil
}


func (db *AeroDB) TakeSeat(tripID int, passenger string, seat int) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	row := db.db.QueryRow("SELECT id FROM Passenger WHERE name=?", passenger)
	var passID int
	if row.Scan(&passID) == sql.ErrNoRows {
		return ErrNotFound;
	}

	row = db.db.QueryRow("SELECT plane_id FROM Trip WHERE id=?", tripID)
	var pID int
	if row.Scan(&pID) == sql.ErrNoRows {
		return ErrNotFound;
	}

	row = db.db.QueryRow("SELECT * FROM Taken WHERE trip_id=? AND place=?", tripID, seat)
	if row.Scan() != sql.ErrNoRows {
		return ErrAlreadyTaken;
	}

	row = db.db.QueryRow("SELECT seats FROM Plane WHERE id=?", pID)
	var seats int
	if row.Scan(&seats) == sql.ErrNoRows {
		return ErrNotFound;
	}

	if (seats < seat) {
		return ErrSeatRange;
	}

	
	_, err := db.db.Exec("INSERT INTO Taken(trip_id, passenger_id, place) VALUES (?, ?, ?)", tripID, passID, seat)

	if (err != nil) {
		return ErrDB
	}
	return nil
}

func (db *AeroDB) GetAllTrips() ([]Trip, error) {
	if (!db.opened) {
		return []Trip{}, ErrNotOpened
	}

	rows, err := db.db.Query("SELECT * FROM Trip")
	if (err == sql.ErrNoRows) {
		return []Trip{}, ErrEmpty
	} else if (err != nil) {
		return []Trip{}, ErrDB
	}
	
	defer rows.Close()
	res := []Trip{}
	for rows.Next() {
		var (
			id int
			company, plane int
    		timeOut, timeIn time.Time
    		townOut, townIn string
		)
		
		if err = rows.Scan(&id, &company, &plane, &timeOut, &timeIn, &townOut, &townIn); err != nil {
			return []Trip{}, ErrDB
		}
		
		new := Trip{id, company, plane, timeOut, timeIn, townOut, townIn}
		res = append(res, new)
	}
	return res, nil
}

func (db *AeroDB) GetTrips(from, to string) ([]Trip, error) {
	if (!db.opened) {
		return []Trip{}, ErrNotOpened
	}

	rows, err := db.db.Query("SELECT * FROM Trip WHERE town_out=? AND town_in=?", from, to)
	if (err == sql.ErrNoRows) {
		return []Trip{}, ErrEmpty
	} else if (err != nil) {
		return []Trip{}, ErrDB
	}
	
	defer rows.Close()
	res := []Trip{}
	for rows.Next() {
		var (
			id int
			company, plane int
    		timeOut, timeIn time.Time
    		townOut, townIn string
		)
		
		if err = rows.Scan(&id, &company, &plane, &timeOut, &timeIn, &townOut, &townIn); err != nil {
			return []Trip{}, ErrDB
		}
		
		new := Trip{id, company, plane, timeOut, timeIn, townOut, townIn}
		res = append(res, new)
	}
	return res, nil
}

func (db *AeroDB) PlanTrip(trip Trip) (TripID int, err error) {
	if (!db.opened) {
		return 0, ErrNotOpened
	}
	if (trip.timeOut.Compare(trip.timeOut) >= 0) {
		return 0, ErrIncorectTime
	}

	_, err = db.db.Exec(`INSERT INTO Trip(company_id, plane_id, time_out, time_in, town_out, town_in)
	            VALUES (?, ?, ?, ?, ?, ?)`, trip.company, trip.plane, trip.timeOut, trip.timeIn, trip.townOut, trip.townIn)

	if (err != nil) {
		return 0, ErrDB
	}

	row := db.db.QueryRow(`SELECT id FROM Trip
	                      WHERE company_id=? AND plane_id=? AND time_out=? AND time_in=? AND town_out=?
						   AND town_in=?`, trip.company, trip.plane, trip.timeOut, trip.timeIn, trip.townOut, trip.townIn)
	if row.Scan(&TripID) != nil {
		return 0, ErrDB
	}
	return TripID, nil
}

func (db *AeroDB) EndTrip(tripID int) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	row := db.db.QueryRow("SELECT * FROM Trip WHERE id=?", tripID)
	if row.Scan() == sql.ErrNoRows {
		return ErrNotFound;
	}

	_, err := db.db.Exec("DELETE FROM Taken WHERE trip_id=?", tripID)
	if (err != nil) {
		return ErrDB
	}


	_, err = db.db.Exec("DELETE FROM Trip WHERE id=?", tripID)
	if (err != nil) {
		return ErrDB
	}
	return nil
}

func (db *AeroDB) DelPlane(name string) (error) {
	if (!db.opened) {
		return ErrNotOpened
	}

	row := db.db.QueryRow("SELECT id FROM Plane WHERE name=?", name)
	var pID int
	if row.Scan(&pID) == sql.ErrNoRows {
		return ErrNotFound;
	}

	rows, err := db.db.Query("SELECT id FROM Trip WHERE plane_id=?", pID)
	if (err != nil && err != sql.ErrNoRows) {
		return ErrDB
	}


	trips := []int{}
	for rows.Next() {
		var tripID int

		if err = rows.Scan(&tripID); err != nil {
			rows.Close()
			return ErrDB
		}

		trips = append(trips, tripID)
	}
	rows.Close()

	for _, id := range trips {
		err = db.EndTrip(id)
		if (err != nil)	{
			return err
		}
	}

	_, err = db.db.Exec("DELETE FROM Plane WHERE name=?", name)
	if (err != nil) {
		return ErrDB
	}
	return nil
}