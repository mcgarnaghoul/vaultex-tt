package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"vaultex-tt/data"
)

var (
	flgHost       string
	flgPort       int
	flgImportFile string
	flgDbFile     string

	con *gorm.DB
)

func init() {
	flag.StringVar(&flgHost, "host", "", "set the hostname to run the server with")
	flag.IntVar(&flgPort, "port", 8080, "set the port to bind to")
	flag.StringVar(&flgImportFile, "import", "", "set the XLSX import filepath")
	flag.StringVar(&flgDbFile, "db", "test.db", "set the SQLite DB file path")
}

func main() {
	fmt.Println("Vaultex - SQL Application")
	flag.Parse()

	// Set up DB connection.
	fmt.Printf("Connecting to %s...\n", flgDbFile)
	var err error
	con, err = gorm.Open(sqlite.Open(flgDbFile), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("unable to connect to DB: %s", err.Error()))
	}

	// Map the schema.
	fmt.Printf("Mapping schema...\n")
	if err = data.AutoMigrate(con); err != nil {
		panic(fmt.Sprintf("error migrating schema: %s", err.Error()))
	}

	// Do any data import needed.
	if flgImportFile != "" {
		fmt.Printf("Importing data from %s...\n", flgImportFile)
		if err := data.ImportFromXLSX(flgImportFile, con); err != nil {
			panic(fmt.Sprintf("error during data import: %s", err.Error()))
		}
	}

	// Register HTTP handlers.
	staticContentHandler := http.FileServer(http.Dir("./static"))
	http.Handle("/", staticContentHandler) // Handle static content on the root.
	http.HandleFunc("/api/", apiHandler)

	// Start HTTP server.
	fmt.Printf("Listening at %s:%d...\n", flgHost, flgPort)
	if err := http.ListenAndServe(getServerAddress(), nil); err != nil {
		fmt.Printf("Fatal error: %s.", err.Error())
	}
}

// Concat host and port.
func getServerAddress() string {
	if flgPort < 0 || flgPort > 65535 {
		panic(fmt.Sprintf("invalid port: %v", flgPort))
	}
	return fmt.Sprintf("%s:%d", flgHost, flgPort)
}

// Handler func for the API route.
func apiHandler(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		code = http.StatusOK
	)
	defer func() {
		// check error
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error occurred: %s.", err.Error())))
		}
		w.WriteHeader(code)
	}()

	// Check the method. Only GET implemented for this but could easily add POST/DELETE.
	switch r.Method {
	case http.MethodGet:
		// Handle get records.
		switch objectName := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/")); objectName {
		case "employees":
			var records []data.Employee
			result := con.Scopes(data.Paginate(r)).Find(&records)

			if err = result.Error; err != nil {
				code = http.StatusInternalServerError
				return
			}

			enc := json.NewEncoder(w)
			if err = enc.Encode(records); err != nil {
				code = http.StatusInternalServerError
			}

		case "organisations":
			var records []data.Organisation
			result := con.Scopes(data.Paginate(r)).Find(&records)

			if err = result.Error; err != nil {
				code = http.StatusInternalServerError
				return
			}

			enc := json.NewEncoder(w)
			if err = enc.Encode(records); err != nil {
				code = http.StatusInternalServerError
			}

		default:
			// TODO plug in logging here
			err = fmt.Errorf("unkown object type: %s", objectName)
			code = http.StatusBadRequest
		}

	default:
		code = http.StatusMethodNotAllowed
	}
}
