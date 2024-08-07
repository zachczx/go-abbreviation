package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go-abbreviation/metaphone3"
	"go-abbreviation/models"

	"github.com/antzucaro/matchr"
	_ "modernc.org/sqlite"
)

type Db struct {
	db *sql.DB
}

const (
	file     string = "./abbreviations.db"
	jsonFile string = "./static/data.json"
)

var db *sql.DB

func DbBase() {
	// Get a database handler
	var err error
	db, err = sql.Open("sqlite", file)
	if err != nil {
		log.Fatal("Error encountered opening DB: ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("Error with pinging DB :", pingErr)
	}
	fmt.Println("Connected!")
}

// Originally tried:
// Jaro-Winkler - doesnt work well comparing a short keyword to long string.
// Aho Corasick - useful for needles in haystack, but not for this
// Soundex - returns for the first word, so not useful
//
// Metaphone3 is roughly about half the time of Jaro-Winkler

func DbRetrieveKeywordMetaphone(keyword string) []models.AbvDb {
	timeStart := time.Now()

	DbBase()

	var prefixQuery string

	if keyword == "" {
		prefixQuery = "SELECT * FROM abv ORDER BY UPPER(short)"
	} else {
		prefixQuery = "SELECT * FROM abv WHERE UPPER(short) LIKE UPPER(?) OR UPPER(long) LIKE UPPER(?) OR metaphone LIKE ? ORDER BY UPPER(short)"
	}

	// Note 1: Without these arguments, Center will match Data Science & Artificial Intelligence. Probably cos it
	// 		   truncated and matches when without max length. So max length is impt.
	// Note 2: If exact match is turned on, center doesnt even match centre, so this is useless.
	// Note 3: Thus need to turn on vowels and max length
	e := &metaphone3.Encoder{EncodeVowels: true, EncodeExact: true, MaxLength: 255}
	keywordMetaScore, _ := e.Encode(keyword)

	// adding this >4 min length control, otherwise 1-3 digit words match everything.
	// But I still want this to use wildcard/contain to match things within words
	if len(keyword) > 4 {
		keyword = "%" + keyword + "%"
		keywordMetaScore = "%" + keywordMetaScore + "%"
	}
	rows, err := db.Query(prefixQuery, keyword, keyword, keywordMetaScore)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// A slice to hold data from returned rows.
	var abbreviations []models.AbvDb

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var item models.AbvDb
		if err := rows.Scan(&item.Short, &item.Long, &item.Initial, &item.Metaphone); err != nil {
			log.Fatal(err)
		}
		abbreviations = append(abbreviations, item)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	timeEnd := time.Since(timeStart)
	fmt.Println("Completed in: ", timeEnd)

	return abbreviations
}

// Using Jaro-Winkler because it performs better than metaphone3
// However trade off is the need to do 2 queries on DB
//
// Once to get all Long, loop thru each to get JW distance, get anything >0.8,
// then use these to get the actual list of objects

func DbRetrieveKeywordJaroWinker(keyword string) []models.AbvDb {
	timeStart := time.Now()
	DbBase()

	if keyword == "" {
		rows, err := db.Query("SELECT * FROM abv;")
		if err != nil {
			fmt.Println("Error with querying: ", err)
		}
		defer rows.Close()

		var abbreviations []models.AbvDb

		for rows.Next() {
			var item models.AbvDb
			if err := rows.Scan(&item.Short, &item.Long, &item.Initial, &item.Metaphone); err != nil {
				log.Fatal(err)
			}
			abbreviations = append(abbreviations, item)
		}
		timeEnd := time.Since(timeStart)
		fmt.Println("Completed in: ", timeEnd)
		return abbreviations
	}

	prepRows, err := db.Query("SELECT long FROM abv;")
	if err != nil {
		fmt.Println("Error with querying: ", err)
	}
	defer prepRows.Close()
	var abvstringMatchedSlice []string

	for prepRows.Next() {
		var prepItem models.AbvDb

		if err := prepRows.Scan(&prepItem.Long); err != nil {
			log.Fatal("testRows scan error: ", err)
		}
		testTempArr := strings.Split(prepItem.Long, " ")
		for _, v := range testTempArr {
			if strings.EqualFold(keyword, v) {
				abvstringMatchedSlice = append(abvstringMatchedSlice, prepItem.Long)
				break
			}

			testDistance := matchr.JaroWinkler(strings.ToUpper(keyword), strings.ToUpper(v), false)

			// Choosing 0.92 >
			// - 0.8 gave a lot of false positives.
			// - 0.9 Prize matched Prime
			if testDistance > 0.92 {
				// Print testDistance
				// fmt.Println(v, testDistance, testItem.Long)
				abvstringMatchedSlice = append(abvstringMatchedSlice, prepItem.Long)
				break
			}
		}

		if err = prepRows.Err(); err != nil {
			log.Fatal("testRows fatal error: ", err)
		}
	}

	// Build the query
	// To Do: Find a way to make this safer, this doesnt feel safe
	var queryBuilder string = `SELECT * FROM abv WHERE UPPER(short) LIKE UPPER(?) OR long IN (`
	for i := 0; i < len(abvstringMatchedSlice); i++ {
		// fmt.Println("Entered the loop", i, " time")
		// SQL escapes using 2 ''. If this isn't escaped, there'll be a
		// SQL logic error: near "s": syntax error (1)
		abvstringMatchedSlice[i] = strings.ReplaceAll(abvstringMatchedSlice[i], "'", "''")

		if i == len(abvstringMatchedSlice)-1 {
			queryBuilder = queryBuilder + `'` + abvstringMatchedSlice[i] + `'`
		} else {
			queryBuilder = queryBuilder + `'` + abvstringMatchedSlice[i] + `',`
		}
	}
	queryBuilder = queryBuilder + `) ORDER BY short;`
	// fmt.Println("Query builder is", queryBuilder)

	// Mutate keyword to add %?%
	// Turning this on when I add an exact match option
	// keyword = "%" + keyword + "%"

	// Execute query
	rows, err := db.Query(queryBuilder, keyword)
	// Need to do this before Close(), else I'll get a
	// http: panic serving [::1]:55879: runtime error: invalid memory address or nil pointer dereference
	// If err != nil then res == nil, so rows panics? i guess.
	if err != nil {
		fmt.Println("SQL query error: ")
		fmt.Println(err)
	}
	defer rows.Close()

	var abbreviations []models.AbvDb

	for rows.Next() {
		var item models.AbvDb
		if err := rows.Scan(&item.Short, &item.Long, &item.Initial, &item.Metaphone); err != nil {
			log.Fatal(err)
		}
		abbreviations = append(abbreviations, item)
	}
	timeEnd := time.Since(timeStart)
	fmt.Println("Completed in: ", timeEnd)
	return abbreviations
}

func DbRetrieveAlphabet(alphabet string) []models.AbvDb {
	timeStart := time.Now()

	DbBase()
	var prefixQuery string
	if alphabet == "" {
		prefixQuery = "SELECT * FROM abv ORDER BY short"
	} else {
		prefixQuery = "SELECT * FROM abv WHERE initial=? ORDER BY short"
	}
	rows, err := db.Query(prefixQuery, alphabet)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var abbreviations []models.AbvDb

	for rows.Next() {
		var item models.AbvDb
		if err := rows.Scan(&item.Short, &item.Long, &item.Initial, &item.Metaphone); err != nil {
			log.Fatal(err)
		}
		abbreviations = append(abbreviations, item)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	timeEnd := time.Since(timeStart)
	fmt.Println("Completed in: ", timeEnd)

	return abbreviations
}

func DbPopulateFromJson() {
	timeStart := time.Now()

	DbBase()
	jsonFileData, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Successfully opened json file")
	defer jsonFileData.Close()

	// Decode json file
	var abvData models.List
	json.NewDecoder(jsonFileData).Decode(&abvData)

	// Prep sql string
	var sqlStr string = "INSERT INTO abv (short, long, initial, metaphone) VALUES "
	vals := []interface{}{}

	// Doing this so that I can implement 0-9 alphabet view
	numberValues := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	e := &metaphone3.Encoder{EncodeVowels: true, EncodeExact: true, MaxLength: 255}
	for i, row := range abvData.List {
		// leaving this in as I'm tired of adding this everytime I want to try something small
		if i < 100000 {
			sqlStr += "(?, ?, UPPER(?), ?),"
			tempVal := string(row.Short[0])
			if contains(numberValues, tempVal) {
				tempVal = "0-9"
			}
			row.Initial = tempVal
			row.Metaphone, _ = e.Encode(row.Long)
			vals = append(vals, row.Short, row.Long, row.Initial, row.Metaphone)
		}
	}
	fmt.Println(vals)
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	fmt.Println(sqlStr)

	// Manually created index
	// Eventually settled on
	// 		CREATE TABLE abv (short VARCHAR(25) collate nocase, long TEXT collate nocase, initial VARCHAR(3) collate nocase, metaphone VARCHAR(255) collate nocase);
	// 		CREATE INDEX short_idx on abv (short COLLATE NOCASE);
	// 		CREATE INDEX initial_idx on abv (initial COLLATE NOCASE);
	// 		DROP INDEX short_upper_idx
	// Left it in, but from earlier testing (not latest indexes) there's no noticeable perf bonus for search queries,
	// e.g. SELECT * FROM abv ORDER BY UPPER(short) take 0.001s anw
	//
	// But for alphabet/initial fetching, helps. Up to 36% improvement (1.5188ms -> 976.9µs) for 'A'.

	result, err := db.Exec(sqlStr, vals...)
	if err != nil {
		fmt.Println("Error executing SQL str: ", err)
	}

	fmt.Println("Populated successfully!")
	// fmt.Println(result)
	_ = result // avoid the declared but not used

	timeEnd := time.Since(timeStart)
	fmt.Println("Completed in: ", timeEnd)
}
