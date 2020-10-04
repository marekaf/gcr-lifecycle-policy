package worker

import (
	"log"
	"sort"
	"strconv"
)

// ByTimeCreated implements sort.Interface for []Digest based on
// the TimeCreatedMs field.
type ByTimeCreated []Digest

func (a ByTimeCreated) Len() int      { return len(a) }
func (a ByTimeCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimeCreated) Less(i, j int) bool {
	inum, err := strconv.Atoi(a[i].TimeCreatedMs)
	if err != nil {
		log.Fatal(err)
	}

	jnum, err := strconv.Atoi(a[j].TimeCreatedMs)
	if err != nil {
		log.Fatal(err)
	}

	return inum < jnum

}

func toSortedSlice(m map[string]Digest) []Digest {

	digests := make([]Digest, 0, len(m))

	for k, v := range m {
		v.Name = k
		digests = append(digests, v)

	}
	sort.Sort(ByTimeCreated(digests))

	return digests
}
