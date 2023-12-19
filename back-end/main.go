package main

import (
	"SocialNetwork/cmd"
)

const port = ":8080"

func main() {
	/*models.ConnectDB()*/
	cmd.Server(port)
}
