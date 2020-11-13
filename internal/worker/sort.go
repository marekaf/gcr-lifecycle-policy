package worker

import (
	"log"
	"sort"
	"strconv"
	"time"
)

// ByTimeCreated implements sort.Interface for []Digest based on
// the TimeCreatedMs field.
type ByTimeCreated []Digest

func (a ByTimeCreated) Len() int      { return len(a) }
func (a ByTimeCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less interprets as i "is newer than" / "happened after" j
func (a ByTimeCreated) Less(i, j int) bool {
	inum, err := strconv.Atoi(a[i].TimeCreatedMs)
	if err != nil {
		log.Fatal(err)
	}

	jnum, err := strconv.Atoi(a[j].TimeCreatedMs)
	if err != nil {
		log.Fatal(err)
	}

	return inum > jnum
}

// ByTimeUploaded implements sort.Interface for []Digest based on
// the TimeUploadedMs field.
type ByTimeUploaded []Digest

func (a ByTimeUploaded) Len() int      { return len(a) }
func (a ByTimeUploaded) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less interprets as i "is newer than" / "happened after" j
func (a ByTimeUploaded) Less(i, j int) bool {

	inum, err := strconv.Atoi(a[i].TimeUploadedMs)
	if err != nil {
		log.Fatal(err)
	}

	jnum, err := strconv.Atoi(a[j].TimeUploadedMs)
	if err != nil {
		log.Fatal(err)
	}

	return inum > jnum
}

func toSortedSlice(sortBy string, m map[string]Digest) []Digest {

	digests := make([]Digest, 0, len(m))

	for k, v := range m {
		v.Name = k
		digests = append(digests, v)
	}

	switch sortBy {

	case "timeCreatedMs":
		sort.Sort(ByTimeCreated(digests))

	case "timeUploadedMs":
		sort.Sort(ByTimeUploaded(digests))

	default:
		log.Fatalf("wrong parameter for sorting. Supported: [timeCreatedMs, timeUploadedMs]. Provided: [%s]", sortBy)
	}

	return digests
}

func olderThanRetention(sortBy string, d Digest, retention time.Time) bool {

	var err error
	var timestamp int64

	switch sortBy {

	case "timeCreatedMs":
		timestamp, err = strconv.ParseInt(d.TimeCreatedMs, 10, 64)

	case "timeUploadedMs":
		timestamp, err = strconv.ParseInt(d.TimeUploadedMs, 10, 64)

	default:
		log.Fatalf("wrong parameter for sorting. Supported: [timeCreatedMs, timeUploadedMs]. Provided: [%s]", sortBy)
	}

	if err != nil {
		log.Fatal(err)
	}

	// retention is in seconds, convert it to ms
	return time.Unix(timestamp/1000, 0).Before(retention)
}
