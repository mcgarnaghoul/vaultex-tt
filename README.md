# Vaultex Technical Test Application
## Overview
Self-contained HTTP server written in golang that serves static web content and API requests to retrieve data from an SQLite DB and display it in a table.  It can also import XLSX data using the `-import` flag.

## Prerequisites
* Install [Go 1.15.5](https://go.dev/doc/install) or higher.
* Open a terminal and grab the dependencies using `go get`:
```shell
go get -u github.com/qax-os/excelize
go get -u gorm.io/gorm
go get -u gorm.io/drivers/sqlite
```
## Building the Application
Once golang is installed, clone the repository, open a terminal in the project directory and run the following
```shell
go build "vaultex-tt.go" && "./vaultex-tt"
```

You should see the following output:
```shell
Vaultex - SQL Application
Connecting to test.db...
Mapping schema...
Listening at :8080...
```

And that's it. You can now view the page in your web browser by going to [http://localhost:8080/](http://localhost:8080/).

![UI screenshot 1](https://raw.githubusercontent.com/mcgarnaghoul/vaultex-tt/main/doc/ui_1.PNG "UI 1")
![UI screenshot 2](https://raw.githubusercontent.com/mcgarnaghoul/vaultex-tt/main/doc/ui_2.PNG "UI 2")

The test data should already be present in test.db, although if it's not...

## Importing Data
Run the application with the `-import` flag to specify an XLSX file to import data from. It is expected that the Excel data fits the schema.

## Other Flags
* `-host` - Specify a hostname to run the server from. Default is blank (uses localhost).
* `-port` - Specify a port to run on. Default is `8080`. Do not include a colon, just the number.
* `-db` - Specify an SQLite database file to connect to. By default the application will try to connect to `test.db` in the working directory.

## Further Improvements
* UI could be moved to a proper JS framework instead of throwing jQuery at it.
* Better error handling.
* Add unit tests.
* Stricter ORM/schema checks.
* Page state not saved after refresh.
* Console error about superfluous WriteHeader call needs fixing.