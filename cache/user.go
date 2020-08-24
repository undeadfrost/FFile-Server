package cache

import (
	cache "FFile-Server/cache/redis"
	db "FFile-Server/db/mysql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username string, password string) bool {
	stmt, err := db.DBConn().Prepare("insert ignore into " +
		"`user`(`username`, `password`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		return true
	}

	return false
}

func LoginUser(username string, password string) bool {
	stmt, err := db.DBConn().Prepare("select password from `user`" +
		" where username=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var hashPassword string
	err = stmt.QueryRow(username).Scan(&hashPassword)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err == nil {
		return true
	}

	return false
}

func SaveSession(username string, sessionToken string, second int) bool {
	redisConn := cache.Pool.Get()
	_, err := redisConn.Do("SETEX", sessionToken, second, username)
	if err == nil {
		return true
	}
	defer redisConn.Close()
	fmt.Println(err.Error())
	return false
}

func AuthSession(sessionToken string) (string, error) {
	redisConn := cache.Pool.Get()
	value, err := redis.String(redisConn.Do("GET", sessionToken))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(value)
	defer redisConn.Close()
	return value, nil
}
