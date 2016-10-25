package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/revel/revel"
	"ringinginross.com/jross/www/app"
	entities "ringinginross.com/jross/www/app/entities"
)


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

// GetGuest returns a guest identified by guestUUID
func GetGuest(guestUUID string) (*entities.Guest, error) {
	app.InitDB()
	query := `
		SELECT
			uuid,
			first_name,
			last_name,
			email,
			party_uuid,
			attending,
			dietary_restriction,
			allergy,
			special_request
		FROM guest
		WHERE uuid = ?`
	revel.INFO.Printf("Query -> %s", query)

	row := app.DB.QueryRow(query, guestUUID)
	var guest entities.Guest
	err := row.Scan(
		&guest.UUID,
		&guest.FirstName,
		&guest.LastName,
		&guest.Email,
		&guest.PartyUUID,
		&guest.Attending,
		&guest.DietaryRestriction,
		&guest.Allergy,
		&guest.SpecialRequest,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			m := fmt.Sprintf("No party uuid found for guest -> Guest: %s", guestUUID)
			return &guest, errors.New(m)
		}
		m := fmt.Sprintf("Database error %s -> Guest: %s", err, guestUUID)
		return &guest, errors.New(m)
	}

	revel.INFO.Printf("Query result -> guestUUID: %s -> %s", *guest.UUID, guest)
	return &guest, nil
}

// GetGuests returns a slice of guests that this guest is part of it's party
func GetGuests(guestUUID string) ([]entities.Guest, error) {
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

	var guests []entities.Guest
	if err == nil {
		for rows.Next() {
			var guest entities.Guest
			rows.Scan(
				&guest.UUID,
				&guest.FirstName,
				&guest.LastName,
				&guest.PartyUUID,
			)
			revel.INFO.Printf("Query response -> Guest name: %s %s", *guest.FirstName, *guest.LastName)
			guests = append(guests, guest)
		}
	}

	return guests, err
}

// GetGuestsByPartyUUID returns a slice of guests for this partyUUID
func GetGuestsByPartyUUID(partyUUID string) ([]entities.Guest, error) {
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

	var guests []entities.Guest
	if err == nil {
		for rows.Next() {
			var guest entities.Guest
			rows.Scan(
				&guest.UUID,
				&guest.FirstName,
				&guest.LastName,
				&guest.PartyUUID,
			)
			revel.INFO.Printf("Query response -> Guest name: %s %s", *guest.FirstName, *guest.LastName)
			guests = append(guests, guest)
		}
	}

	return guests, err
}


func SetGuestInformation(guestUUID, email, dietary, special string, attending, allergy bool) error {
	app.InitDB()

	// verify this guest exists
	_, err := GetGuest(guestUUID)
	if err != nil {
		revel.ERROR.Print(err)
		m := fmt.Sprintf("Error retrieving guest: %s", guestUUID)
		return errors.New(m)
	}

	query := `
		UPDATE
			guest
		SET
			email = ?,
			attending = ?,
			dietary_restriction = ?,
			allergy = ?,
			special_request = ?
		WHERE
			uuid = ?
	`
	revel.INFO.Printf("Query -> %s", query)

	_, err = app.DB.Exec(
		query,
		email,
		attending,
		dietary,
		allergy,
		special,
		guestUUID,
	)

	if err != nil {
		m := fmt.Sprintf("Database error %s -> Guest: %s", err, guestUUID)
		return errors.New(m)
	}
	revel.INFO.Printf("Guest successfully updated -> Guest: %s", guestUUID)
	return nil
}


func SetDeclineMessage(partyUUID, message string) error {
	query := `
		UPDATE
			party
		SET
			message = ?
		WHERE
			uuid = ?
	`
	revel.INFO.Printf("Query -> %s", query)
	_, err := app.DB.Exec(
		query,
		message,
		partyUUID,
	)

	if err != nil {
		m := fmt.Sprintf("Database error %s -> Party: %s", err, partyUUID)
		return errors.New(m)
	}
	revel.INFO.Printf("Party decline updated -> Party: %s", partyUUID)
	return nil
}


func SetConfirmMessage(partyUUID, message string, uber, transportation bool) error {
	query := `
		UPDATE
			party
		SET
			message = ?,
			uber = ?,
			transportation = ?
		WHERE
			uuid = ?
	`
	revel.INFO.Printf("Query -> %s", query)
	_, err := app.DB.Exec(
		query,
		message,
		uber,
		transportation,
		partyUUID,
	)

	if err != nil {
		m := fmt.Sprintf("Database error %s -> Party: %s", err, partyUUID)
		return errors.New(m)
	}
	revel.INFO.Printf("Party confirm updated -> Party: %s", partyUUID)
	return nil
}
