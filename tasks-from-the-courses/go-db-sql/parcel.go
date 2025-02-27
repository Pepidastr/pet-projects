package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:Client, :Status, :Address, :CreatedAt)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("CreatedAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	res := s.db.QueryRow("select * from parcel where number = :number", sql.Named("number", number))
	p := Parcel{}
	// заполните объект Parcel данными из таблицы
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	query, err := s.db.Query("select * from parcel where client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for query.Next() {
		p := Parcel{}
		err = query.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}
	err = query.Err()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Query("update parcel set status = :status where number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Query("update parcel set address = :address where number = :number AND status = :statusRegistered",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("statusRegistered", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {

	_, err := s.db.Query("delete from parcel where number = :number AND status = :statusRegistered",
		sql.Named("number", number),
		sql.Named("statusRegistered", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
