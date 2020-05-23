package main

import (
	"TAX/server"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*import (
	_ "github.com/lib/pq"
)*/
type Person struct {
	Name string
	Age  int
}

func main() {

	err := godotenv.Load(`.env`)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//var dbSource = flag.String("DB_SOURCE", os.Getenv("DB_SOURCE"), "Database url connection string.")
	//var dbType = flag.String("DB_TYPE", os.Getenv("DB_TYPE"), "Database type.")
	var serverPort = flag.String("SERVER_PORT", os.Getenv("SERVER_PORT"), "Port to run server on.")

	flag.Parse()

	/*dbConn, err := driver.ConnectDB(*dbSource, *dbType)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}*/

	server.Run(*serverPort)

}

/*func show()
{
	fmt.Printf("request success")
}*/
