package worker

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestExtractRepositoryFromImage(t *testing.T) {

	tables := []struct {
		x string
		y Repository
	}{
		{"eu.gcr.io/mygcp-project/acme/webapp:master",
			Repository{
				RepositoryPrefix: "eu.gcr.io/mygcp-project/acme/",
				ImageName:        "webapp",
				Tag:              "master",
			},
		},
		{"us.gcr.io/mygcp-project/mycorp/webapp",
			Repository{
				RepositoryPrefix: "us.gcr.io/mygcp-project/mycorp/",
				ImageName:        "webapp",
				Tag:              "",
			},
		},
		{"gcr.io/mygcp-project/mycorp/webapp:latest",
			Repository{
				RepositoryPrefix: "gcr.io/mygcp-project/mycorp/",
				ImageName:        "webapp",
				Tag:              "latest",
			},
		},
	}

	for _, table := range tables {
		repo := extractRepositoryFromImage(table.x)
		diff := cmp.Diff(repo, table.y)
		if diff != "" {
			t.Fatalf("extractRepositoryFromImage mismatch (-want +got):\n%v", diff)
		}
	}
}

func TestDaysToTime(t *testing.T) {

	start := time.Now()

	tables := []struct {
		x int
		y time.Time
	}{
		{10, start.AddDate(0, 0, -10)},
		{0, start},
		{666, start.AddDate(0, 0, -666)},
	}

	for _, table := range tables {
		time := daysToTime(start, table.x)
		diff := cmp.Diff(time, table.y)
		if diff != "" {
			t.Fatalf("daysToTime mismatch (-want +got):\n%v", diff)
		}
	}
}
