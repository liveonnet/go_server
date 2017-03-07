package user_lib

import (
	"testing"
)

var g_ticket string

func Test_Create(t *testing.T) {
	ticket, err := CreateTicket(135792468)
	if err == nil {
		t.Logf("ticket %s", ticket)
		g_ticket = ticket
	} else {
		t.Logf("error: %s", err)
	}
}

func Test_Validate(t *testing.T) {
	//	ticket := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiLvv70iLCJleHAiOjMxNTM2MDAwLCJpc3MiOiJ4eGRiIn0.ue6hdn89dtMg8koH9h-7oF5eCPtNUQLxaMqSTfErzhE"
	uid, err := ValidateTicket(g_ticket)
	if err == nil {
		t.Logf("uid %v", uid)
	} else {
		t.Logf("error: %s", err)
	}
}
