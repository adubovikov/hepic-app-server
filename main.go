package main

import (
	"hepic-app-server/v2/cmd"
)

// @title HEPIC App Server v2 API
// @version 2.0.0
// @description Advanced REST API Server for HEPIC App Server v2
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.hepic.local/support
// @contact.email support@hepic.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT authorization token. Format: "Bearer {token}"

func main() {
	cmd.Execute()
}
