package users

import (
	"testing"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestUserValidate(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		username     string
		firstname    string
		lastname     string
		expectError  bool
	}{
		{
			"Valid User",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			false,
		},
		{
			"Invalid email address",
			"*12djda+1:~",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			true,
		},
		{
			"Password length too short",
			"rioaishii@gmail.com",
			"pswd",
			"pswd",
			"rioishii",
			"Rio",
			"Ishii",
			true,
		},
		{
			"Password and PasswordConf does not match",
			"rioaishii@gmail.com",
			"password1234",
			"grW12309d",
			"rioishii",
			"Rio",
			"Ishii",
			true,
		},
		{
			"Empty user name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"",
			"Rio",
			"Ishii",
			true,
		},
		{
			"Contains spaces in user name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rio ishii",
			"Rio",
			"Ishii",
			true,
		},
	}

	for _, c := range cases {
		nu := NewUser{c.email, c.password, c.passwordConf, c.username, c.firstname, c.lastname}
		err := nu.Validate()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error validating new user: %v", c.name, err)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one", c.name)
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		username     string
		firstname    string
		lastname     string
		expectError  bool
	}{
		{
			"Successful converting NewUser to User",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			false,
		},
		{
			"Invalid new user",
			"rioaishii@gmail.com",
			"awd@q309f",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			true,
		},
	}

	for _, c := range cases {
		nu := NewUser{c.email, c.password, c.passwordConf, c.username, c.firstname, c.lastname}
		_, err := nu.ToUser()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error validating new user: %v", c.name, err)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one", c.name)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		username     string
		firstname    string
		lastname     string
		expected     string
	}{
		{
			"Valid full name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"Rio Ishii",
		},
		{
			"Empty full name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"",
			"",
			"",
		},
		{
			"Empty first name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"",
			"Ishii",
			"Ishii",
		},
		{
			"Empty last name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"",
			"Rio",
		},
	}

	for _, c := range cases {
		nu := NewUser{c.email, c.password, c.passwordConf, c.username, c.firstname, c.lastname}
		user, err := nu.ToUser()
		if err != nil {
			t.Errorf("unexpected error converting NewUser to User")
		}
		if output := user.FullName(); output != c.expected {
			t.Errorf("case %s: expected %s but got %s", c.name, c.expected, output)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		username     string
		firstname    string
		lastname     string
		authPassword string
		expectError  bool
	}{
		{
			"Successful authentication",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"password1234",
			false,
		},
		{
			"Password hash mismatch",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"ad120981209",
			true,
		},
	}

	for _, c := range cases {
		nu := NewUser{c.email, c.password, c.passwordConf, c.username, c.firstname, c.lastname}
		user, err := nu.ToUser()
		if err != nil {
			t.Errorf("unexpected error converting NewUser to User")
		}
		err = user.Authenticate(c.authPassword)
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected authentication error: %v", c.name, err)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one", c.name)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		username     string
		firstname    string
		lastname     string
		updateFirst  string
		updateLast   string
		expectError  bool
		expectedName string
	}{
		{
			"Valid full name update",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"John",
			"Doe",
			false,
			"John Doe",
		},
		{
			"Update first name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"John",
			"",
			false,
			"John Ishii",
		},
		{
			"Update last name",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"",
			"Doe",
			false,
			"Rio Doe",
		},
		{
			"Invalid update",
			"rioaishii@gmail.com",
			"password1234",
			"password1234",
			"rioishii",
			"Rio",
			"Ishii",
			"",
			"",
			true,
			"Rio Ishii",
		},
	}

	for _, c := range cases {
		nu := NewUser{c.email, c.password, c.passwordConf, c.username, c.firstname, c.lastname}
		update := Updates{c.updateFirst, c.updateLast}
		user, err := nu.ToUser()
		if err != nil {
			t.Errorf("unexpected error converting NewUser to User")
		}
		err = user.ApplyUpdates(&update)
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error applying updates: %v", c.name, err)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one", c.name)
		}
		if output := user.FullName(); output != c.expectedName {
			t.Errorf("case %s: expected %s but got %s", c.name, c.expectedName, output)
		}
	}
}
