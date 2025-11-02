package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func connectDB() (*sqlx.DB, error) {
	// connect to the PostgreSQL database
	connStr := "user=user password=password dbname=mydatabase host=localhost port=5430 sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InsertUser(db *sqlx.DB, user User) error {
	// insert a new user into the database
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	// retrieve all users from the database
	var users []User
	err := db.Select(&users, "SELECT * FROM users ORDER BY id")
	return users, err
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	// retrieve a user by ID
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	// transfer balance from one user to another within a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	var sender User
	err = tx.Get(&sender, "SELECT * FROM users WHERE id=$1 FOR UPDATE", fromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("sender not found: %v", err)
	}

	if sender.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id=$2", amount, fromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to withdraw: %v", err)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id=$2", amount, toID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deposit: %v", err)
	}

	return tx.Commit()
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalln("Database connection failed:", err)
	}
	defer db.Close()

	// add some users
	user1 := User{Name: "Daniyar", Email: "daniyar@kbtu.kz", Balance: 100.0}
	user2 := User{Name: "Daniyal", Email: "daniyal@kbtu.kz", Balance: 50.0}

	if err := InsertUser(db, user1); err != nil {
		log.Println("InsertUser error:", err)
	}
	if err := InsertUser(db, user2); err != nil {
		log.Println("InsertUser error:", err)
	}

	// get and display all users
	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("GetAllUsers error:", err)
	} else {
		fmt.Println("\n=== Users ===")
		for _, u := range users {
			fmt.Printf("[%d] %s | %s | balance: %.2f\n", u.ID, u.Name, u.Email, u.Balance)
		}
	}

	// transfer balance
	fmt.Println("\n=== Transfer 30.0 from Daniyar to Daniyal ===")
	if err := TransferBalance(db, 1, 2, 30.0); err != nil {
		log.Println("Transfer error:", err)
	}

	// 4️ display users after transfer
	users, _ = GetAllUsers(db)
	fmt.Println("\n=== After Transfer ===")
	for _, u := range users {
		fmt.Printf("[%d] %s | balance: %.2f\n", u.ID, u.Name, u.Balance)
	}
}
