package worker

import (
	"strconv"
	"testing"
	"time"
)

func Test_olderThanRetention(t *testing.T) {
	type args struct {
		sortBy    string
		d         Digest
		retention time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"created is older than now()", args{"timeCreatedMs", Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{""},
			TimeCreatedMs:  strconv.FormatInt(time.Now().AddDate(0, -1, 0).Unix()*1000, 10),
			TimeUploadedMs: strconv.FormatInt(time.Now().AddDate(0, -1, 0).Unix()*1000, 10),
		}, time.Now()}, true},
		{"uploaded is older than now()", args{"timeUploadedMs", Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{""},
			TimeCreatedMs:  "0",
			TimeUploadedMs: strconv.FormatInt(time.Now().AddDate(0, -1, 0).Unix()*1000, 10),
		}, time.Now()}, true},
		{"created is not older than 3 years", args{"timeCreatedMs", Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{""},
			TimeCreatedMs:  strconv.FormatInt(time.Now().AddDate(-2, -1, 0).Unix()*1000, 10),
			TimeUploadedMs: strconv.FormatInt(time.Now().AddDate(-1, -1, 0).Unix()*1000, 10),
		}, time.Now().AddDate(-3, 0, 0)}, false},
		{"uploaded is not older than 2 months", args{"timeUploadedMs", Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{""},
			TimeCreatedMs:  "0",
			TimeUploadedMs: strconv.FormatInt(time.Now().AddDate(0, -1, 0).Unix()*1000, 10),
		}, time.Now().AddDate(0, -2, 0)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := olderThanRetention(tt.args.sortBy, tt.args.d, tt.args.retention); got != tt.want {
				t.Errorf("olderThanRetention() = %v, want %v", got, tt.want)
			}
		})
	}
}
