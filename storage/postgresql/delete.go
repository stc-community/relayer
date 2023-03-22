package postgresql

import "time"

func (b *Postgres) DeleteEvent(id string, pubkey string) error {
	_, err := b.DB.Exec("DELETE FROM event WHERE id = $1 AND pubkey = $2", id, pubkey)
	return err
}

func (b *Postgres) Clean() {
	for {
		time.Sleep(60 * time.Minute)
		b.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
	}
}
