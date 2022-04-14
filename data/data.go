package data

import (
	"fmt"
	"github.com/qax-os/excelize"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

const (
	maxRecordCount = 1000
)

// Grab schema information and populate ORM models.
func AutoMigrate(con *gorm.DB) error {
	if err := con.AutoMigrate(&Organisation{}); err != nil {
		return err
	}

	if err := con.AutoMigrate(&Employee{}); err != nil {
		return err
	}

	return nil
}

// Returns a pointer to a scoped DB for pagination.
func Paginate(r *http.Request) func(*gorm.DB) *gorm.DB {
	return func(con *gorm.DB) *gorm.DB {
		if err := r.ParseForm(); err != nil {
			con.AddError(err)
			return con
		}

		pageParam := r.Form.Get("p")
		page := 1
		if pageParam != "" {
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				con.AddError(err)
				return con
			}
		}
		if page < 1 {
			page = 1
		}

		countParam := r.Form.Get("c")
		count := 100
		if countParam != "" {
			var err error
			count, err = strconv.Atoi(countParam)
			if err != nil {
				con.AddError(err)
				return con
			}
		}
		if count < 1 || count > maxRecordCount {
			count = maxRecordCount
		}

		offset := (page - 1) * count
		return con.Offset(offset).Limit(count)
	}
}

// Import data from an XLSX file.
func ImportFromXLSX(path string, con *gorm.DB) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Iterate through each sheet and build records as map[string]interface{}.
	sheets := f.GetSheetList()
	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}

		var columnNames []string
		switch strings.ToLower(sheetName) {
		case "employee":
			records := make([]Employee, len(rows)-1)
			for i, row := range rows {
				if i == 0 {
					columnNames = row
					continue
				}

				values := make(map[string]interface{})
				for j, col := range columnNames {
					if j >= len(columnNames) {
						break
					}
					values[col] = row[j]
				}
				records[i-1] = Employee{
					OrganisationNumber: values["OrganisationNumber"].(string),
					FirstName:          values["FirstName"].(string),
					LastName:           values["LastName"].(string),
				}
			}

			// Batch insert
			if err = con.Create(&records).Error; err != nil {
				return err
			}

		case "organisation":
			records := make([]Organisation, len(rows)-1)
			for i, row := range rows {
				if i == 0 {
					columnNames = row
					continue
				}

				values := make(map[string]interface{})
				for j, col := range columnNames {
					if j >= len(columnNames) {
						break
					}
					values[col] = row[j]
				}
				records[i-1] = Organisation{
					OrganisationName:   values["OrganisationName"].(string),
					OrganisationNumber: values["OrganisationNumber"].(string),
					AddressLine1:       values["AddressLine1"].(string),
					AddressLine2:       values["AddressLine2"].(string),
					AddressLine3:       values["AddressLine3"].(string),
					AddressLine4:       values["AddressLine4"].(string),
					Town:               values["Town"].(string),
					Postcode:           values["Postcode"].(string),
				}
			}

			// Batch insert
			if err = con.Create(&records).Error; err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown object type from sheet name: %s", sheetName)
		}
	}

	return nil
}
