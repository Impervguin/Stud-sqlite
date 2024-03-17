package aerodb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if (err != nil) {
		return "", err
	}
	return string(b), nil
}

func diffSql(pathDB1, pathDB2 string) (string, error) {
	cmd := exec.Command("sqldiff", pathDB1, pathDB2)
	b, err := cmd.Output()
	return string(b), err
}

func createDatabase(fname string) (error) {
	cmd := exec.Command("python3", "./scripts/createaero.py", fname)
	b, err := cmd.Output()
	if (string(b) != "") {
		return fmt.Errorf(string(b))
	}
	if (err != nil) {
		return err
	}
	return nil
}

func createTestDataBases() (string, string, error) {
	tbase, err := os.CreateTemp("./", "tmp")
	if (err != nil) {
		return "", "", err
	}

	tmod, err := os.CreateTemp("./", "tmp")
	if (err != nil) {
		os.Remove(tbase.Name())
		return "", "", err
	}

	err = createDatabase(tbase.Name())
	if (err != nil) {
		os.Remove(tbase.Name())
		os.Remove(tmod.Name())
		return "", "", err
	}

	err = createDatabase(tmod.Name())
	if (err != nil) {
		os.Remove(tbase.Name())
		os.Remove(tmod.Name())
		return "", "", err
	}
	return tbase.Name(), tmod.Name(), err
}

func readTest(dir string) (out, diff string, err error) {
	b, err := os.ReadFile(dir + "eout")
	if (err != nil) {
		return "", "", err
	}
	out = string(b)

	b, err = os.ReadFile(dir + "ediff.sql")
	if (err != nil) {
		return "", "", err
	}
	diff = string(b)
	return out, diff, nil
}

func errMessage(err error) (string) {
	if err == nil {
		return "nil"
	} else {
		return err.Error()
	}
}

// Positive Test 1: AddPassenger
func TestAddPassengerPositive(t *testing.T) {
	dir := "tests/pos1/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Error(err)
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Error(err)
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.AddPassenger("Mark")

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 2: AddCompany
func TestAddCompanyPositive(t *testing.T) {
	dir := "tests/pos2/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.AddCompany("StudAirlines")

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 3: AddPlane
func TestAddPlanePositive(t *testing.T) {
	dir := "tests/pos3/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.AddPlane("Antosha", "S7", 192)

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 4: DelPlane
func TestDelPlanePositive(t *testing.T) {
	dir := "tests/pos4/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.DelPlane("AirBus A310")

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 5: DelCompany
func TestDelCompanyPositive(t *testing.T) {
	dir := "tests/pos5/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.DelCompany("S7", "Red Wings")

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 6: TakeSeat
func TestTakeSeatPositive(t *testing.T) {
	dir := "tests/pos6/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.TakeSeat(1, "Batman", 10)

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 7: EndTrip
func TestEndTripPositive(t *testing.T) {
	dir := "tests/pos7/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	funcErr := db.EndTrip(1)

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 8: CreateTrip
func TestCreateTripPositive(t *testing.T) {
	dir := "tests/pos8/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	trip := CreateTrip(0, 1, 3, time.Unix(1707642000, 0), time.Unix(1707661680, 0), "Moscow", "Tokyo")
	_, funcErr := db.PlanTrip(trip)

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) 
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 9: GetFreeSeats
func TestGetFreeSeatsPositive(t *testing.T) {
	dir := "tests/pos9/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	seats, funcErr := db.GetFreeSeats(4)

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) + "\n"
	out += strings.Trim(strings.Join(strings.Fields(fmt.Sprint(seats)), ", "), "[]")
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}

// Positive test 10: GetTrips
func TestGetTrips(t *testing.T) {
	dir := "tests/pos10/" // Директория с данными для теста
	tbase, tmod, err := createTestDataBases() // Создаёт два временных файла и записываем в них
	// тестовую базу данных
	if (err != nil) {
		t.Errorf("Error while creating temp files for test: %v", err.Error())
		return
	}
	// Удаляем временные файлы
	defer os.Remove(tbase)
	defer os.Remove(tmod)

	// Считываем тестовые данные
	// eout - возвращаемые значения функции
	// ediff - отличие базы данных после выполнения операции от стандартной
	eout, ediff, err := readTest(dir)
	if (err != nil) {
		t.Errorf("Error while fetching test data: %v", err.Error())
		return
	}
	eout = strings.Trim(eout, "\n")

	// Начало теста
	db := AeroDB{}
	err = db.OpenDB(tmod)
	if (err != nil) {
		t.Errorf("Cannot open database: %v", err.Error())
		return
	}
	
	// Тестовое действие
	arr, funcErr := db.GetTrips("Moscow", "New-york")
	s := ""
	for _, el := range arr {
		s += "\n" + fmt.Sprint(el)
	}

	err = db.CloseDB()
	if (err != nil) {
		t.Errorf("Cannot close database: %v", err.Error())
		return
	}

	// Получение вывода в строковом формате
	out := errMessage(funcErr) + s
	// Сравнение полученной базы данных исходной
	diff, err := diffSql(tbase, tmod)
	
	if (err != nil) {
		t.Errorf("Cannot compare databases: %v", err.Error())
		return
	}
	// Анализ полученных результатов с ожидаемыми
	// Сравниваем возвращаемые значения функции
	if (out != eout) {
		t.Errorf("Incorrect output\nGot:\n%v\nExpected:\n%v", out, eout)
	}
	// Сравниваем изменение базы данных
	if (diff != ediff) {
		t.Errorf("Incorrect database action\nGot:\n%v\nExpected:\n%v", diff, ediff)
	}
}