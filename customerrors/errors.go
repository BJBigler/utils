package customerrors

import (
	"fmt"
)

//UserUnregistered ...
type UserUnregistered struct {
	UserID  string
	Message string
}

func (e *UserUnregistered) Error() string {
	return fmt.Sprintf(e.Message)
}

//-----------

//DatabaseMissing ...
type DatabaseMissing struct {
	Message string
}

func (e *DatabaseMissing) Error() string {
	return fmt.Sprintf(e.Message)
}

//-----------

//UserUnauthorized ...
type UserUnauthorized struct {
	UserID  string
	Message string
}

//Error ...
func (e *UserUnauthorized) Error() string {

	if e.Message == "" && e.UserID == "" {
		return "user is unauthorized to view resource"
	}

	if e.UserID != "" && e.Message == "" {
		return fmt.Sprintf("user (%s) is unauthorized to view resource", e.UserID)
	}

	return fmt.Sprintf("user is unauthorized to view resource: %s", e.Message)
}

//PrettyPrint outputs error message for user
func (e *UserUnauthorized) PrettyPrint() string {
	return fmt.Sprintf("Your user permission level prevents you from performing this action or accessing the requested resource. (%s)", e.UserID)
}

//UserNotFound ...
type UserNotFound struct {
	Message string
}

//Error ...
func (e *UserNotFound) Error() string {
	return "user was not found in the database"
}

//PrettyPrint outputs error message for user
func (e *UserNotFound) PrettyPrint() string {
	return `We could not find your user record in the database. 
			You may have signed in with different credentials than when you signed up, 
			or with a different authentication provider. If you believe you should 
			have access but haven't yet created a profile, please contact an administrator
			for a signup link. You will need to <a href="/logout">logout</a> to return to 
			anonymous browsing.`
}

//-----------

//NotUniqueResult ...
type NotUniqueResult struct {
	Message   string
	ResultCnt int
}

//NotUniqueResult resource returned more than one result
//when it should have returned one
func (e *NotUniqueResult) Error() string {
	return e.Message
}

//PrettyPrint ...
func (e *NotUniqueResult) PrettyPrint() string {
	return fmt.Sprintf("The query performed should have returned exactly one result, but instead return %v.", e.ResultCnt)
}

//-----------

//NoResult ...
type NoResult struct {
	Message string
	Err     error
}

func (e *NoResult) Error() string {
	return e.Message
}

//------------

//NoProfile is for user that have successfully logged into the authenticator,
//but do not have a profile on the system.
type NoProfile struct {
	UserID      string
	Institution string
	Err         error
	Message     string
}

func (e *NoProfile) Error() string {
	return e.Message
}
