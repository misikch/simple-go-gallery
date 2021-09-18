package main

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestStorageAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	st := NewStorage(db)

	userID := 1
	path := "testpath"
	testPhoto := &Photo{
		UserID: userID,
		Path: path,
	}

	//ok query
	mock.
		ExpectExec("INSERT INTO `photos`").
		WithArgs(userID, path).
		WillReturnResult(sqlmock.NewResult(1,1))

	err = st.Add(testPhoto)
	if err != nil {
		t.Errorf("unexpected err on method Add: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil { //не было ли лишних запросов? Только то, что ожидали?
		t.Errorf("where were unexpected sqlmock error: %v", err)
	}

	//sql.DB returned error check
	mock.
		ExpectExec("INSERT INTO `photos`").
		WithArgs(userID, path).
		WillReturnError(fmt.Errorf("bad query"))

	err = st.Add(testPhoto)
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil { //не было ли лишних запросов? Только то, что ожидали?
		t.Errorf("where were unexpected sqlmock error: %v", err)
	}

	//lastInsertId error check
	mock.
		ExpectExec("INSERT INTO `photos`").
		WithArgs(userID, path).
		WillReturnResult(sqlmock.NewResult(0,0))

	err = st.Add(testPhoto)
	if err == nil {
		t.Errorf("expected err about wrong 'insertId', got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil { //не было ли лишних запросов? Только то, что ожидали?
		t.Errorf("where were unexpected sqlmock error: %v", err)
	}
}

func TestStorageGetPhotos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	st := NewStorage(db)

	//good query
	testUserId := 777
	rows := sqlmock.NewRows([]string{"id", "user_id", "path"})
	expect := []*Photo{
		{
			ID: 1,
			UserID: testUserId,
			Path: "path_1",
		},
		{
			ID: 2,
			UserID: testUserId,
			Path: "path_2",
		},
	}
	for _, p := range expect {
		rows.AddRow(p.ID, p.UserID, p.Path)
	}

	mock.
		ExpectQuery("SELECT `id`, `user_id`, `path` FROM `photos` WHERE user_id").
		WithArgs(testUserId).
		WillReturnRows(rows)

	items, err := st.GetPhotos(testUserId)
	if err != nil {
		t.Errorf("unexpected getPhotos() error: %v", err)
		return
	}

	if !reflect.DeepEqual(expect, items) {
		t.Errorf("results not match")
	}

	//row scan error
	rows = sqlmock.NewRows([]string{"id", "user_id"})
	for _, p := range expect {
		rows.AddRow(p.ID, p.UserID)
	}

	mock.
		ExpectQuery("SELECT `id`, `user_id`, `path` FROM `photos` WHERE user_id").
		WithArgs(testUserId).
		WillReturnRows(rows)

	_, err = st.GetPhotos(testUserId)
	if err == nil {
		t.Errorf("expected error while scanning rows to struct, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil { //не было ли лишних запросов? Только то, что ожидали?
		t.Errorf("where were unexpected sqlmock error: %v", err)
	}
}