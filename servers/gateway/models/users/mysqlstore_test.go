package users

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	nu := NewUser{
		Email:        "rioaishii@gmail.com",
		Password:     "password1234",
		PasswordConf: "password1234",
		UserName:     "rioishii",
		FirstName:    "Rio",
		LastName:     "Ishii",
	}

	err = nu.Validate()
	if err != nil {
		t.Fatalf("cannot validate new user: %v", err)
	}

	user, err := nu.ToUser()
	if err != nil {
		t.Fatalf("cannot convert to user: %v", err)
	}

	expectedSQL := regexp.QuoteMeta(sqlInsertUser)

	var newID int64 = 1

	mock.ExpectExec(expectedSQL).
		WithArgs(
			user.Email,
			user.PassHash,
			user.UserName,
			user.FirstName,
			user.LastName,
			user.PhotoURL,
		).
		WillReturnResult(sqlmock.NewResult(newID, 1))

	insertedUser, err := sqlStore.Insert(user)
	if err != nil {
		t.Fatalf("unexpected error during successful insert: %v", err)
	}
	if insertedUser == nil {
		t.Fatal("nil user returned from insert")
	} else if insertedUser.ID != newID {
		t.Fatalf("incorrect new ID: expected %d but got %d", newID, insertedUser.ID)
	}
}

func TestUserInsertFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	nu := NewUser{
		Email:        "rioaishii@gmail.com",
		Password:     "password1234",
		PasswordConf: "password1234",
		UserName:     "rioishii",
		FirstName:    "Rio",
		LastName:     "Ishii",
	}

	err = nu.Validate()
	if err != nil {
		t.Fatalf("cannot validate new user: %v", err)
	}

	user, err := nu.ToUser()
	if err != nil {
		t.Fatalf("cannot convert to user: %v", err)
	}

	expectedSQL := regexp.QuoteMeta(sqlInsertUser)

	mock.ExpectExec(expectedSQL).
		WithArgs(
			user.Email,
			user.PassHash,
			user.FirstName,
			user.LastName,
			user.PhotoURL,
		).
		WillReturnError(fmt.Errorf("some error"))

	_, err = sqlStore.Insert(user)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}).
		AddRow(1, "test@gmail.com", "test", "username", "first", "last", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByID)

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(userMockRows)

	user, err := sqlStore.GetByID(1)
	expectedUser := User{ID: 1, Email: "test@gmail.com", PassHash: user.PassHash, UserName: "username", FirstName: "first", LastName: "last", PhotoURL: "testtest"}
	if err != nil {
		t.Fatalf("unexpected error during successful select: %v", err)
	}
	if user == nil {
		t.Fatal("nil user returned from select")
	} else if !reflect.DeepEqual(user, &expectedUser) {
		t.Fatalf("incorrect user: expected %v but got %v", expectedUser, user)
	}
}

func TestGetUserByIDScanFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "first_name", "last_name", "photo_url"}).
		AddRow(1, "test@gmail.com", "test", "first", "last", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByID)

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(userMockRows)

	_, err = sqlStore.GetByID(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)
	id := int64(1)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}).
		AddRow(id, "rioaishii@gmail.com", "password", "rioishii", "rio", "ishii", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByEmail)
	mock.ExpectQuery(expectedSQL).
		WithArgs("rioaishii@gmail.com").
		WillReturnRows(userMockRows)

	user, err := sqlStore.GetByEmail("rioaishii@gmail.com")
	expectedUser := User{ID: 1, Email: "rioaishii@gmail.com", PassHash: user.PassHash, UserName: "rioishii", FirstName: "rio", LastName: "ishii", PhotoURL: "testtest"}
	if err != nil {
		t.Fatalf("unexpected error during successful select: %v", err)
	}
	if user == nil {
		t.Fatal("nil user returned from select")
	} else if !reflect.DeepEqual(user, &expectedUser) {
		t.Fatalf("incorrect user: expected %v but got %v", expectedUser, user)
	}
}

func TestGetUserByEmailFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	expectedSQL := regexp.QuoteMeta(sqlGetUserByEmail)
	mock.ExpectQuery(expectedSQL).
		WithArgs("test@gmail.com").
		WillReturnError(fmt.Errorf("email not found"))

	_, err = sqlStore.GetByEmail("test@gmail.com")
	if err == nil {
		t.Fatalf("unexpected error, found invalid user %v", err)
	}
}

func TestGetUserByEmailScanFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "first_name", "last_name", "photo_url"}).
		AddRow(1, "test@gmail.com", "test", "first", "last", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByEmail)

	mock.ExpectQuery(expectedSQL).
		WithArgs("test@gmail.com").
		WillReturnRows(userMockRows)

	_, err = sqlStore.GetByEmail("test@gmail.com")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetUserByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)
	id := int64(1)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}).
		AddRow(id, "rioaishii@gmail.com", "password", "rioishii", "rio", "ishii", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByUserName)
	mock.ExpectQuery(expectedSQL).
		WithArgs("rioishii").
		WillReturnRows(userMockRows)

	user, err := sqlStore.GetByUserName("rioishii")
	if err != nil {
		t.Fatalf("unexpected error during successful select: %v", err)
	}
	expectedUser := User{ID: 1, Email: "rioaishii@gmail.com", PassHash: user.PassHash, UserName: "rioishii", FirstName: "rio", LastName: "ishii", PhotoURL: "testtest"}
	if user == nil {
		t.Fatal("nil user returned from select")
	} else if !reflect.DeepEqual(user, &expectedUser) {
		t.Fatalf("incorrect user: expected %v but got %v", expectedUser, user)
	}
}

func TestGetUserByUserNameFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	expectedSQL := regexp.QuoteMeta(sqlGetUserByUserName)
	mock.ExpectQuery(expectedSQL).
		WithArgs("test").
		WillReturnError(fmt.Errorf("username not found"))

	_, err = sqlStore.GetByUserName("test")
	if err == nil {
		t.Fatalf("unexpected error, found invalid user %v", err)
	}
}

func TestGetUserByUserNameScanFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "first_name", "last_name", "photo_url"}).
		AddRow(1, "test@gmail.com", "test", "first", "last", "testtest")

	expectedSQL := regexp.QuoteMeta(sqlGetUserByUserName)

	mock.ExpectQuery(expectedSQL).
		WithArgs("test").
		WillReturnRows(userMockRows)

	_, err = sqlStore.GetByUserName("test")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)

	updateID := int64(1)
	updateUser := Updates{FirstName: "John", LastName: "Doe"}

	expectedRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}).
		AddRow(1, "rioaishii@gmail.com", "password", "rioishii", "John", "Doe", "testtest")

	expectedSQLUpdate := regexp.QuoteMeta(sqlUpdateUser)
	expectedSQLGet := regexp.QuoteMeta(sqlGetUserByID)

	mock.ExpectExec(expectedSQLUpdate).
		WithArgs(
			updateUser.FirstName,
			updateUser.LastName,
			updateID,
		).
		WillReturnResult(sqlmock.NewResult(0, updateID))
	mock.ExpectQuery(expectedSQLGet).
		WithArgs(updateID).
		WillReturnRows(expectedRows)

	user, err := sqlStore.Update(updateID, &updateUser)

	expectedUser := User{ID: 1, Email: "rioaishii@gmail.com", PassHash: user.PassHash, UserName: "rioishii", FirstName: "John", LastName: "Doe", PhotoURL: "testtest"}
	if err != nil {
		t.Fatalf("unexpected error during successful update: %v", err)
	} else if !reflect.DeepEqual(user, &expectedUser) {
		t.Fatalf("incorrect user: expected %v but got %v", expectedUser, user)
	}
}

func TestUpdateUserFailRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()
	sqlStore := NewSQLStore(db)
	updateID := int64(1)
	updateUser := Updates{FirstName: "John", LastName: "Doe"}
	userMockRows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}).
		AddRow(updateID, "rioaishii@gmail.com", "password", "rioishii", "rio", "ishii", "testtest")
	expectedSQLUpdate := regexp.QuoteMeta(sqlUpdateUser)
	expectedSQLGet := regexp.QuoteMeta(sqlGetUserByID)
	mock.ExpectExec(expectedSQLUpdate).
		WithArgs(
			updateUser.FirstName,
			updateID,
		).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectQuery(expectedSQLGet).
		WithArgs(updateID).
		WillReturnRows(userMockRows)
	_, err = sqlStore.Update(updateID, &updateUser)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)
	deleteID := int64(2)

	expectedSQLDelete := regexp.QuoteMeta(sqlDeleteUser)
	expectSQLGet := regexp.QuoteMeta(sqlGetUserByID)

	mock.ExpectExec(expectedSQLDelete).
		WithArgs(deleteID).
		WillReturnResult(sqlmock.NewResult(0, deleteID))
	mock.ExpectQuery(expectSQLGet).
		WithArgs(deleteID).
		WillReturnError(fmt.Errorf("user unfound"))

	err = sqlStore.Delete(deleteID)
	if err != nil {
		t.Fatalf("unexpected error during successful update: %v", err)
	}
	_, err = sqlStore.GetByID(deleteID)
	if err == nil {
		t.Fatalf("unexpected error returned deleted user %v", err)
	}
}

func TestDeleteUserFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	sqlStore := NewSQLStore(db)
	deleteID := int64(5000)

	expectedSQLDelete := regexp.QuoteMeta(sqlDeleteUser)
	expectSQLGet := regexp.QuoteMeta(sqlGetUserByID)

	mock.ExpectExec(expectedSQLDelete).
		WithArgs(deleteID).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectQuery(expectSQLGet).
		WithArgs(deleteID).
		WillReturnError(fmt.Errorf("user unfound"))

	err = sqlStore.Delete(deleteID)
	if err == nil {
		t.Fatalf("Expected error")
	}
}
