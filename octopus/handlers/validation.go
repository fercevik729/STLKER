package handlers

import (
	"regexp"
	"strings"
)

// validatePortfolio validates a portfolio's name
func validatePortfolio(port *Portfolio) (bool, string) {
	// Check the length of the name and if it contains spaces
	if len(port.Name) < 3 || len(port.Name) > 30 || strings.Contains(port.Name, " ") {
		return false, "Name must be between 3 and 30 characters and cannot contain spaces"
	}
	// Check if the name is alphanumeric
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)
	matches := re.FindAllString(port.Name, -1)

	return len(matches) == 1, "Name must be alphanumeric"

}

// validateUser validates a user
func validateUser(usr User) (bool, string) {
	// Check the lengths
	if len(usr.Username) < 6 || len(usr.Username) > 30 || len(usr.Password) < 10 || len(usr.Password) > 100 {
		return false, "Username must be between 6 and 30 characters. Password must be between 10 and 100 characters"
	}
	// Check if the username or pwd contain invalid chars or the password contains the username
	if strings.ContainsAny(usr.Username, "(){}[]|!%^@:;&_'-+<>") || strings.ContainsAny(usr.Password, "(){}[]|!%^@:;&_'-+<>") || strings.Contains(usr.Password, usr.Username) {
		return false, "Username or password contains invalid characters"
	}
	return true, ""
}
