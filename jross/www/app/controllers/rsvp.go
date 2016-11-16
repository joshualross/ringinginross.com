package controllers

import (
	"database/sql"
	"math/rand"

	"github.com/revel/revel"
	db "ringinginross.com/jross/www/app/repositories"
)

// Rsvp Controller
type Rsvp struct {
	*revel.Controller
}

var invalidAccessError = []string{
	"Really excited to come to the wedding, huh?  Sorry, that wasn't correct though",
}

var userNotFoundError = []string{
	"That name wasn't found in the registry, did you make a typo?",
	"Who do you think you are!? No, really, who are you? I couldn't find you on the list",
	"You forgot how to spell your name. For reference, check your invite for the correct spelling ;)",
}

var databaseError = []string{
	"Uh oh, that wasn't supposed to happen!  Do you mind reporting this to Josh or Eunice?",
	"Whoopsie daisy!  Looks like Josh isn't as good at programming as he thought...",
}
var cookieError = []string{
	"Something happened that wasn't supposed to happen.  Sorry, enter your name again please :)",
}

// Index render template requesting password entry
func (c Rsvp) Index() revel.Result {
	return c.Redirect(Rsvp.Name)
	// if c.Session["access"] != "" {
	// 	return c.Redirect(Rsvp.Name)
	// }
	// if c.Session["user"] != "" {
	// 	return c.Redirect(Rsvp.Detail)
	// }
	// return c.Render()
}

// Access checks for correct password
func (c Rsvp) Access(password string) revel.Result {
	revel.INFO.Print("Access action")
	// if password != "foo" {
	// 	c.Flash.Error(getErrorMessage(invalidAccessError))
	// 	return c.Redirect(Rsvp.Index)
	// }
	// c.Session["access"] = "1"

	return c.Redirect(Rsvp.Name)
}

// Name renders the name submission
func (c Rsvp) Name() revel.Result {
	return c.Render()
}

// NameSubmit queries for given first and last name
func (c Rsvp) NameSubmit(response, firstName, lastName string) revel.Result {

	revel.INFO.Printf("Received parameters -> Response: %s, First: %s, Last: %s", response, firstName, lastName)

	guestUUID, err := db.GetGuestUUID(firstName, lastName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Flash.Error(getErrorMessage(userNotFoundError))
		} else {
			c.Flash.Error(getErrorMessage(databaseError))
		}
		return c.Redirect(Rsvp.Name)
	}
	c.Session["uuid"] = guestUUID

	// Save the response
	partyUUID, _ := db.GetPartyUUID(guestUUID)
	_ = db.SetGuestResponse(partyUUID, response == "Yes")

	if response != "Yes" {
		// Decline the invitation
		return c.Redirect(Rsvp.Decline)
	}
	return c.Redirect(Rsvp.Detail)
}

// Decline verifies the user is coming
func (c Rsvp) Decline() revel.Result {
	guestUUID := c.Session["uuid"]
	partyUUID, err := db.GetPartyUUID(guestUUID)
	if err != nil {
		// render without party uuid, we will try again in the submit
		partyUUID = ""
	} else {
		_ = db.SetGuestResponse(partyUUID, false)
	}
	return c.Render(partyUUID)
}

// DeclineSubmit saves the guests response
func (c Rsvp) DeclineSubmit() revel.Result {
	// TODO save the decline message, if not empty
	guestUUID := c.Session["uuid"]
	message := c.Request.PostForm.Get("message")
	partyUUID, err := db.GetPartyUUID(guestUUID)
	if err != nil {
		revel.ERROR.Printf("Unable to save message for guest %s: %s", guestUUID, message)
	}

	if err := db.SetDeclineMessage(partyUUID, message); err != nil {
		revel.ERROR.Printf("Database error: %s", err)
		revel.ERROR.Printf("Unable to save message for party %s: %s", partyUUID, message)
	}

	return c.Render()
}

// Detail asks for specifics of the rsvp
func (c Rsvp) Detail() revel.Result {
	guestUUID := c.Session["uuid"]
	if guestUUID == "" {
		c.Flash.Error(getErrorMessage(cookieError))
		return c.Redirect(Rsvp.Name)
	}

	guests, err := db.GetGuests(guestUUID)
	if err != nil {
		c.Flash.Error(getErrorMessage(databaseError))
		return c.Redirect(Rsvp.Name)
	}

	partyUUID := guests[0].PartyUUID
	return c.Render(guests, partyUUID)
}

// DetailSubmit saves the guests response
func (c Rsvp) DetailSubmit(partyUUID string) revel.Result {
	revel.INFO.Print("Detail Submit for Party %s %T", partyUUID, partyUUID)
	partyUUID = c.Request.PostForm.Get("partyUUID")

	guests, err := db.GetGuestsByPartyUUID(partyUUID)
	if err != nil {
		c.Flash.Error(getErrorMessage(databaseError))
		return c.Redirect(Rsvp.Detail)
	}

	for _, guest := range guests {
		email := c.Request.PostForm.Get(*guest.UUID + ".Email")
		dietary := c.Request.PostForm.Get(*guest.UUID + ".DietaryRestriction")
		attending := c.Request.PostForm.Get(*guest.UUID + ".Attending") == "on"
		allergy := c.Request.PostForm.Get(*guest.UUID + ".Allergy") == "on"
		special := c.Request.PostForm.Get(*guest.UUID + ".SpecialRequest")


		revel.INFO.Printf(
			"Updating guest %s: (%s, %s, %v, %v, %s)",
			*guest.UUID,
			email,
			dietary,
			attending,
			allergy,
			special,
		)
		err := db.SetGuestInformation(*guest.UUID, email, dietary, special, attending, allergy)
		if err != nil {
			revel.ERROR.Printf("Update error %s", err)
		}
	}

	return c.Redirect(Rsvp.Confirm)
}

// Confirm shows a thank you message
func (c Rsvp) Confirm() revel.Result {
	guestUUID := c.Session["uuid"]
	partyUUID, err := db.GetPartyUUID(guestUUID)
	if err != nil {
		// render without party uuid, we will try again in the submit
		partyUUID = ""
	}
	return c.Render(partyUUID)
}


// ConfirmSubmit stores the last submitted data
func (c Rsvp) ConfirmSubmit() revel.Result {
	guestUUID := c.Session["uuid"]
	message := c.Request.PostForm.Get("message")
	uber := c.Request.PostForm.Get("uber") == "on"
	transportation := c.Request.PostForm.Get("transportation") == "on"
	partyUUID, err := db.GetPartyUUID(guestUUID)
	if err != nil {
		revel.ERROR.Printf("Unable to save message for guest %s: %s", guestUUID, message)
	}

	if err := db.SetConfirmMessage(partyUUID, message, uber, transportation); err != nil {
		revel.ERROR.Printf("Database error: %s", err)
		revel.ERROR.Printf("Unable to save message for party %s: %s", partyUUID, message)
	}

	return c.Render()
}

func getErrorMessage(m []string) string {
	return m[rand.Int()%len(m)]
}
