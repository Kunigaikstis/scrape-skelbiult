package listing

import "database/sql"

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (r *Storage) GetById(id int) (Ad, error) {
	row := r.db.QueryRow("SELECT id, url, created, price, sq_footage, street, neighbourhood, location FROM listing WHERE id = ?", id)

	listing := Ad{}

	err := row.Scan(&listing.Id, &listing.Url, &listing.Created, &listing.Price, &listing.SqFootage, &listing.Street, &listing.Neighbourhood, &listing.Location)

	return listing, err
}

func (r *Storage) Add(listing Ad) error {
	stmt, err := r.db.Prepare("INSERT INTO listing(id, url, price, sq_footage, street, neighbourhood, location) values(?,?,?,?,?,?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(listing.Id, listing.Url, listing.Price, listing.SqFootage, listing.Street, listing.Neighbourhood, listing.Location)

	if err != nil {
		return err
	}

	return nil
}
