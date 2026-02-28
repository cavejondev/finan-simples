// Package person é o pacote de pessoa no domain
package person

import "time"

// Person representa a entidade de pessoa no dominio
type Person struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
