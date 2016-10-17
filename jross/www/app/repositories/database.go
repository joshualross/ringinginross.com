package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/revel/revel"
	"ringinginross.com/jross/www/app"
)

// Guest struct
type Guest struct {
	UUID               string
	FirstName          string
	LastName           string
	Email              string
	PartyUUID          string
	DietaryRestriction string
	Allergy            bool
	SpecialRequest     string
	Attending          bool
}

// GetGuestUUID performs a db lookup and returns a guest UUID
func GetGuestUUID(firstName, lastName string) (string, error) {
	app.InitDB()

	query := "SELECT uuid FROM guest WHERE first_name LIKE ? AND last_name LIKE ?"
	revel.INFO.Printf("Query -> %s", query)
	row := app.DB.QueryRow(query, firstName, lastName)

	var guestUUID string
	err := row.Scan(&guestUUID)
	if err != nil {
		revel.ERROR.Printf("Query error -> %s", err)
		return "", err
	}

	revel.INFO.Printf("Query result -> guestUUID: %s", guestUUID)
	return guestUUID, nil
}

// SetGuestResponse saves an RSVP response for a guest
func SetGuestResponse(partyUUID string, response bool) error {
	app.InitDB()

	query := "UPDATE guest SET attending = ? WHERE party_uuid = ?"
	revel.INFO.Printf("Query -> %s", query)

	result, err := app.DB.Exec(query, response, partyUUID)
	revel.INFO.Printf("%s -> %s", result, err)
	return nil
}

// GetPartyUUID performs a lookup for this guest's party UUID and returns it
func GetPartyUUID(guestUUID string) (string, error) {
	app.InitDB()
	query := "SELECT party_uuid FROM guest WHERE uuid = ?"
	revel.INFO.Printf("Query -> %s", query)

	row := app.DB.QueryRow(query, guestUUID)
	var partyUUID string
	err := row.Scan(&partyUUID)
	if err != nil {
		revel.ERROR.Printf("Query error -> %s", err)
		if err == sql.ErrNoRows {
			m := fmt.Sprintf("No party uuid found for guest -> Guest: %s", guestUUID)
			return "", errors.New(m)
		}
		m := fmt.Sprintf("Database error -> Guest: %s", guestUUID)
		return "", errors.New(m)
	}

	revel.INFO.Printf("Query result -> partyUUID: %s", partyUUID)
	return partyUUID, nil
}

// GetGuests returns a slice of guests that this guest is part of it's party
func GetGuests(guestUUID string) ([]Guest, error) {
	app.InitDB()
	query := `
		SELECT
			uuid,
			first_name,
			last_name,
			party_uuid
		FROM guest
		WHERE party_uuid IN (
			SELECT party_uuid FROM guest WHERE uuid = ?
		)`
	revel.INFO.Printf("Query -> %s", query)

	rows, err := app.DB.Query(query, guestUUID)
	defer rows.Close()

	var guests []Guest
	if err == nil {
		for rows.Next() {
			var guest Guest
			rows.Scan(
				&guest.UUID,
				&guest.FirstName,
				&guest.LastName,
				&guest.PartyUUID,
			)
			revel.INFO.Printf("Query response -> Guest name: %s %s", guest.FirstName, guest.LastName)
			guests = append(guests, guest)
		}
	}

	return guests, err
}

// GetGuestsByPartyUUID returns a slice of guests for this partyUUID
func GetGuestsByPartyUUID(partyUUID string) ([]Guest, error) {
	app.InitDB()
	query := `
		SELECT
			uuid,
			first_name,
			last_name,
			party_uuid
		FROM guest
		WHERE party_uuid = ?`
	revel.INFO.Printf("Query -> %s", query)

	rows, err := app.DB.Query(query, partyUUID)
	defer rows.Close()

	var guests []Guest
	if err == nil {
		for rows.Next() {
			var guest Guest
			rows.Scan(
				&guest.UUID,
				&guest.FirstName,
				&guest.LastName,
				&guest.PartyUUID,
			)
			revel.INFO.Printf("Query response -> Guest name: %s %s", guest.FirstName, guest.LastName)
			guests = append(guests, guest)
		}
	}

	return guests, err
}
