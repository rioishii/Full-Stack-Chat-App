package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/indexes"
)

// SQLStore keeps tracks of the current active database connection so that we don't need to open a new connection
// every time we execute a query
type SQLStore struct {
	db *sql.DB
}

//NewSQLStore constructs a new MySQLStore
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

var v struct {
	Data []User `json:"data"`
}

const sqlGetAllUsers = "select id, user_name, first_name, last_name from users"
const sqlColumnListWithID = "id, email, pass_hash, user_name, first_name, last_name, photo_url"
const sqlColumnListNoID = "email, pass_hash, user_name, first_name, last_name, photo_url"
const sqlGetUserByID = "select " + sqlColumnListWithID + " from users where id = ?"
const sqlGetUserByEmail = "select " + sqlColumnListWithID + " from users where email = ?"
const sqlGetUserByUserName = "select " + sqlColumnListWithID + " from users where user_name = ?"
const sqlInsertUser = "insert into users(" + sqlColumnListNoID + ") values (?,?,?,?,?,?)"
const sqlUpdateUser = "update users set first_name = ?, last_name = ? where id = ?"
const sqlDeleteUser = "delete from users where id = ?"

//GetByID returns the User with the given ID
func (ms *SQLStore) GetByID(id int64) (*User, error) {
	rows, err := ms.db.Query(sqlGetUserByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting next row: %v", err)
	}
	return &user, nil
}

//GetByEmail returns the User with the given email
func (ms *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err := ms.db.Query(sqlGetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting next row: %v", err)
	}
	return &user, nil
}

//GetByUserName returns the User with the given Username
func (ms *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err := ms.db.Query(sqlGetUserByUserName, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting next row: %v", err)
	}
	return &user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *SQLStore) Insert(user *User) (*User, error) {
	result, err := ms.db.Exec(sqlInsertUser, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("error inserting new row: %v", err)
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting new ID: %v", err)
	}

	user.ID = newID
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	_, err := ms.db.Exec(sqlUpdateUser, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("error updating row: %v", err)
	}
	user, err := ms.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting updated user: %v", err)
	}
	return user, nil
}

//Delete deletes the user with the given ID
func (ms *SQLStore) Delete(id int64) error {
	_, err := ms.db.Exec(sqlDeleteUser, id)
	if err != nil {
		return fmt.Errorf("error deleting row: %v", err)
	}
	return nil
}

//GetAllUsers gets all the rows in the user table and inserts them
//into the trie
func (ms *SQLStore) GetAllUsers() (*indexes.Trie, error) {
	trie := indexes.NewTrie()
	rows, err := ms.db.Query(sqlGetAllUsers)
	if err != nil {
		fmt.Print("no users stored in db")
		return trie, nil
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		firstname := strings.ToLower(user.FirstName)
		fsplit := strings.Split(firstname, " ")
		for i := range fsplit {
			fsplit[i] = strings.TrimSpace(fsplit[i])
			trie.Add(fsplit[i], user.ID)
		}

		lastname := strings.ToLower(user.LastName)
		lsplit := strings.Split(lastname, " ")
		for i := range lsplit {
			lsplit[i] = strings.TrimSpace(lsplit[i])
			trie.Add(lsplit[i], user.ID)
		}

		trie.Add(strings.ToLower(user.UserName), user.ID)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error getting next row: %v", err)
	}
	return trie, nil
}
