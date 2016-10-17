package controllers

import (
	"database/sql"
	"math/rand"

	"github.com/revel/revel"
	"ringinginross.com/jross/www/app"
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
	defer app.DB.Close()

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

	return c.Render()
}

// DeclineSubmit saves the guests response
func (c Rsvp) DeclineSubmit() revel.Result {

	// TODO mark all in party yes or no? -> happens in NameSubmit
	return c.Render()
}

// Detail asks for specifics of the rsvp
func (c Rsvp) Detail() revel.Result {
	defer app.DB.Close()
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
	defer app.DB.Close()
	revel.INFO.Print("Detail Submit for Party %s %T", partyUUID, partyUUID)
	revel.INFO.Print("")
	revel.INFO.Printf("params: %s %T", c.Params, c.Params)
    if err := c.Request.ParseForm(); err != nil {
        // handle error
    }
	for key, values := range c.Request.PostForm {
	    revel.INFO.Printf("Key: %s %T -> value: %s %T", key, key, values, values)
	}
	partyUUID = c.Request.PostForm.Get("partyUUID")
	revel.INFO.Printf("partyUUID: %s %T", partyUUID, partyUUID)

	guests, err := db.GetGuestsByPartyUUID(partyUUID)
	if err != nil {
		c.Flash.Error(getErrorMessage(databaseError))
		return c.Redirect(Rsvp.Detail)
	}

	for _, guest := range guests {
		guest.Email = c.Request.PostForm.Get(guest.UUID + ".Email")
		guest.DietaryRestriction = c.Request.PostForm.Get(guest.UUID + ".DietaryRestriction")
		guest.Attending = c.Request.PostForm.Get(guest.UUID + ".Attending") == "on"
		guest.Allergy = c.Request.PostForm.Get(guest.UUID + ".Allergy") == "on"
		guest.SpecialRequest = c.Request.PostForm.Get(guest.UUID + ".SpecialRequest")
		revel.INFO.Printf("guest %s: %s", guest.UUID, guest)
	}

	//revel.INFO.Printf("form: %s %T", c.Request.PostForm, c.Request.PostForm)
	//uuids := c.Request.PostForm.Get("uuid[]")
	//revel.INFO.Printf("uuids: %s %T", uuids, uuids)
	//uuids = c.Request.PostForm.Get("uuid[]")
	//revel.INFO.Printf("uuids: %s %T", uuids, uuids)
	//revel.INFO.Print("")
    //uuids := c.Params.Form.Get("uuid")
    //revel.INFO.Printf("values: %s %T", uuids)
    //revel.INFO.Print("")
    //for i := range uuids {
    //    revel.INFO.Println(uuids[i])
    //}
	//guests := c.Params.Form
	//revel.INFO.Printf("%s", UUID)
	//revel.INFO.Printf("%s", email)
	//revel.INFO.Printf("%s", specialRequest)
	//revel.INFO.Printf("%s", dietaryRestrictions)
	//revel.INFO.Printf("%s", c.Params.Form)
	//guests := c.Params.Form.Get("guests")
	//revel.INFO.Printf("%s", c.Params.Form.Get("guests"))
	//revel.INFO.Printf("%s", guests)
	//revel.INFO.Printf("%s", c.Params.Values.Get("UUID"))
	//revel.INFO.Printf("%s", c.Params.Values)
	return c.Redirect(Rsvp.Confirm)
}

// Confirm shows a thank you message
func (c Rsvp) Confirm() revel.Result {
	// TODO save song selection
	return c.Render()
}


// ConfirmSubmit stores the last submitted data
func (c Rsvp) ConfirmSubmit() revel.Result {

	return c.Render()
}

func getErrorMessage(m []string) string {
	return m[rand.Int()%len(m)]
}
