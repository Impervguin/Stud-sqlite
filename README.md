# Домашнее задание №5

Домашнее задание предназначено для изучения базы данных SQL через систему управления sqlite3 на языке go.

# Задание

Необходимо написать библиотеку для работы с базой данных перелётов, которая имеет следующую структуру:

![Структура базы данных](/readme/aero.png)

 - **Plane** - Хранит информацию о самолётах
 - **Passenger** - Хранит информацию о всех пассажирах
 - **Company** - Хранит информацию о всех самолётах
 - **Trip** - Хранит информацию о запланированных перелётах
 - **Taken** - Хранит информацию о занятых местах в самолёте в определённой поездке

Необходимо реализовать функции открытия и закрытия структуры **aero_db** и методы этой структуры, так чтобы они соответствовали интерфейсу:

```go
type sqlite3_db interface {
    plan_trip(company, plane string, time_out, time_in time.Time, town_out, town_in string) (error)
    end_trip(trip_id int, name string, seat int) (error)
    get_trips(from, to string) ([]Trip, error)
    take_seat(trip_id int, name string, int seat) (error)
    get_free_seats(trip_id int) ([]int, error)
    add_company(name string) (error)
    del_company(name string) (error)
    add_plane(name string) (error)
    del_plane(name string) (error)
    add_passenger(name, surname, middle_name string) (error)
}

type Trip struct {
    id int
    company, plane int
    time_out, time_in time.Time
    town_out, town_in int
}

func open_db(fname string) (*sqlite3_db, error) {}
func close_db(*sqlite3_db) (error) {}
```

# Описание методов и функций

#### Функция `open_db`

`Вход:` Название файла базы данных.
`Выход:` Указатель на структуру базы данных, ошибка(или nil)

Функция подготавливает структуры базы данных, открывая файл sqlite3 и подготавливает структуры к работе(если в вашей структуре есть ещё что-то помимо объекта базы)

#### Функция `close_db`

`Вход:` Указатель на структуру базы данных
`Выход:` Ошибка(или nil)

Функция завершает работу с базой данных и закрывает файл с ней, затем удаляеи структуру

***to be continued***