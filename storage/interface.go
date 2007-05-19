package storage

type Storer interface {
	GetConnection()
}
"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"