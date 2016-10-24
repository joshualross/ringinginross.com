package entities

import "fmt"


// Guest struct
type Guest struct {
	UUID               *string
	FirstName          *string
	LastName           *string
	Email              *string
	PartyUUID          *string
	DietaryRestriction *string
	Allergy            *bool
	SpecialRequest     *string
	Attending          *bool
}

// Update accepts input and updates the struct
func (g Guest) Update(u *Guest) {

	if g.Email != u.Email {
		g.Email = u.Email
	}
	if g.DietaryRestriction != u.DietaryRestriction {
		g.DietaryRestriction = u.DietaryRestriction
	}
	if g.Allergy != u.Allergy {
		g.Allergy = u.Allergy
	}
	if g.Attending != u.Attending {
		g.Attending = u.Attending
	}
	if g.SpecialRequest != u.SpecialRequest {
		g.SpecialRequest = u.SpecialRequest
	}
}

func (g Guest) GetUpdateParams() (email, dietary, special string, attending, allergy bool) {
	if g.Email != nil {
		email = *g.Email
	}
	if g.DietaryRestriction != nil {
		dietary = *g.DietaryRestriction
	}
	if g.SpecialRequest != nil {
		special = *g.SpecialRequest
	}
	if g.Attending != nil {
		attending = *g.Attending
	}
	if g.Allergy != nil {
		allergy = *g.Allergy
	}
	return email, dietary, special, attending, allergy
}


func (g Guest) String() string {
	m := fmt.Sprintf("(%s :: [%s, %s", *g.UUID, *g.FirstName, *g.LastName)
	if g.Email != nil {
		m += fmt.Sprintf(", email:%s", *g.Email)
	} else {
		m += fmt.Sprintf(", email:%v", g.Email)
	}

	if g.Attending != nil {
		m += fmt.Sprintf(", attending:%s", *g.Attending)
	} else {
		m += fmt.Sprintf(", attending:%v", g.Attending)
	}
	if g.DietaryRestriction != nil {
		m += fmt.Sprintf(", dietary:%s", *g.DietaryRestriction)
	} else {
		m += fmt.Sprintf(", dietary:%v", g.DietaryRestriction)
	}
	if g.Allergy != nil {
		m += fmt.Sprintf(", allergy:%s", *g.Allergy)
	} else {
		m += fmt.Sprintf(", allergy:%v", g.Allergy)
	}
	if g.SpecialRequest != nil {
		m += fmt.Sprintf(", special:%s", *g.SpecialRequest)
	} else {
		m += fmt.Sprintf(", special:%v", g.SpecialRequest)
	}
	m += "])"

	return m
}
