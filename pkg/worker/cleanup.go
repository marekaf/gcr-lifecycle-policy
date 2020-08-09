package worker

import (
	"log"
	"strconv"
	"time"
)

func cleanup(list FilteredList) {

}

func olderThanRetention(d Digest, retention time.Time) bool {

	timecreated, err := strconv.ParseInt(d.TimeCreatedMs, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	// retention is in days, convert it to ms
	return time.Unix(timecreated/1000, 0).Before(retention)
}

// HandleCleanup function
func HandleCleanup(c Config) string {

	token := getToken(c.CredsFile)

	catalog := fetchCatalog(token)

	filteredCatalog := filterCatalog(catalog, c.RepoFilter)

	list := fetchTags(token, filteredCatalog)

	cleanupList := filter(list)

	printBeforeCleanup(cleanupList)

	cleanup(cleanupList)

	return "ahoj"
}
