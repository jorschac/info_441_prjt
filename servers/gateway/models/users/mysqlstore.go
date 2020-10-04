package users

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type SQLStore struct {
	//DB
	DB *sql.DB
}


func (sqls *SQLStore) GetByID(id int64) (*User, error) {
	u := &User{}

	insq := "select * from user where id = ?"

	err := sqls.DB.QueryRow(insq, id).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName, &u.FirstName, &u.LastName, &u.PhotoURL, &u.Description)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByEmail returns the User with the given email
func (sqls *SQLStore) GetByEmail(email string) (*User, error) {
	u := &User{}
	insq := "select * from user where email = ?"

	err := sqls.DB.QueryRow(insq, email).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName, &u.FirstName, &u.LastName, &u.PhotoURL, &u.Description)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByUserName returns the User with the given Username
func (sqls *SQLStore) GetByUserName(username string) (*User, error) {
	u := &User{}
	insq := "select * from user where username = ?"

	err := sqls.DB.QueryRow(insq, username).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName, &u.FirstName, &u.LastName, &u.PhotoURL, &u.Description)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (sqls *SQLStore) Insert(user *User) (*User, error) {
	insq := "insert into user(email, passhash, username, first_name, last_name, photourl, description) values (?,?,?,?,?,?,?)"

	res, err := sqls.DB.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL, user.Description)

	if err != nil {
		return nil, err
	} else {
		id, idErr := res.LastInsertId()
		if idErr != nil {
			return nil, idErr
		} else {
			user.ID = id
			return user, nil
		}
	}
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (sqls *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	u := &User{}

	insq := "update user set description = ? where id = ?"

	res, err := sqls.DB.Exec(insq, updates.Description, id)

	if err != nil {
		return nil, err
	} else {
		_, idErr := res.RowsAffected()
		if idErr != nil {
			return nil, idErr
		} else {
			rowErr := sqls.DB.QueryRow("select * from user where id = ?", id).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName, &u.FirstName, &u.LastName, &u.PhotoURL, &u.Description)
			if rowErr != nil {
				return nil, rowErr
			} else {
				return u, nil
			}
		}
	}
}

//Delete deletes the user with the given ID
func (sqls *SQLStore) Delete(id int64) error {

	insq := "delete from user where id = ?"

	_, err := sqls.DB.Exec(insq, id)

	if err != nil {
		return err
	}

	return nil
}

func (sqls *SQLStore) TrackLogin(id int64, loginDate time.Time, ipaddress string) error {
	insq := "insert into logins(id, login_date, ip_address) values (?,?,?)"
	intid := int(id)

	_, err := sqls.DB.Exec(insq, intid, loginDate, ipaddress)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (sqls *SQLStore) Follow(following int64, follower int64) error {
	curIng := ""
	curEr := ""
	errQuery := sqls.DB.QueryRow("select user_following,  user_followed from follow where user_following = ? and user_followed = ?", following, follower).Scan(&curIng, &curEr)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			insq := "insert into follow(user_following, user_followed, date_followed) values(?,?,?)"
			_, err := sqls.DB.Exec(insq, following, follower, time.Now())
			if err != nil {
				return err
			}
		} else {
			return errQuery
		}
	} else {
		insq := "delete from follow where user_following = ? and user_followed = ?"
		_, err := sqls.DB.Exec(insq, following, follower)

		if err != nil {
			return err
		}
	}

	return nil
}
