package aerodb

import (
	"errors"
)


var (
    ErrNotOpened    = errors.New("database is not opened")
    ErrAlreadyIn    = errors.New("element already in database")
    ErrNotFound     = errors.New("element not found")
    ErrEmpty        = errors.New("empty result")
    ErrSeatRange    = errors.New("incorrect seat number") 
    ErrAlreadyTaken = errors.New("seat already taken")   
    ErrFile         = errors.New("cannot open the file")
    ErrDBFormat     = errors.New("not correct format of database")
    ErrIncorectTime = errors.New("incorrect time period")
    
    // На случай неконтролируемых ошибок работы с базой данных
    ErrDB           = errors.New("unknown mistakes with database")
) 