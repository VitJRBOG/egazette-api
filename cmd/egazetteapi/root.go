package egazetteapi

import (
	"egazette-api/internal/server"
)

// Execute starts the main functions of the program.
func Execute() {
	server.Up()
}
