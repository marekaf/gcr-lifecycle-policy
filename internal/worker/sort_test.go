package worker

import (
	"reflect"
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

func Test_toSortedSlice(t *testing.T) {
	type args struct {
		sortBy string
		m      map[string]Digest
	}
	tests := []struct {
		name string
		args args
		want []Digest
	}{
		// one item should be always sorted
		{"one item sort", args{"timeCreatedMs",
			map[string]Digest{
				"sha256:1366ef2a3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263": {
					ImageSizeBytes: "29905102",
					LayerID:        "",
					MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
					Tag:            []string{"master.6d902732dc8c0e19725eaa40c7a860a4c02ef406"},
					TimeCreatedMs:  "1556399434238",
					TimeUploadedMs: "1556399443871",
				},
			},
		}, []Digest{
			{
				ImageSizeBytes: "29905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"master.6d902732dc8c0e19725eaa40c7a860a4c02ef406"},
				TimeCreatedMs:  "1556399434238",
				TimeUploadedMs: "1556399443871",
				Name:           "sha256:1366ef2a3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263",
			},
		},
		},
		// three items
		{"three items sort by timeCreatedMs", args{"timeCreatedMs",
			map[string]Digest{
				"sha256:1366ef2a3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263": {
					ImageSizeBytes: "29905102",
					LayerID:        "",
					MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
					Tag:            []string{"master.6d902732dc8c0e19725eaa40c7a860a4c02ef406"},
					TimeCreatedMs:  "1556399434238",
					TimeUploadedMs: "1556399443871",
				},
				"sha256:2066edba3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263": {
					ImageSizeBytes: "299051042",
					LayerID:        "",
					MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
					Tag:            []string{"stage"},
					TimeCreatedMs:  "1556399434248",
					TimeUploadedMs: "1556399443871",
				},
				"sha256:5e6a2a225050edcb62cea2a01f4f1e4b2610b6c9e98e8b347f78c49fdc05aff7": {
					ImageSizeBytes: "28905102",
					LayerID:        "",
					MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
					Tag:            []string{"dev", "latest"},
					TimeCreatedMs:  "1556399434237",
					TimeUploadedMs: "1556399443871",
				},
			},
		}, []Digest{
			{
				ImageSizeBytes: "299051042",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"stage"},
				TimeCreatedMs:  "1556399434248", // newest by 11ms
				TimeUploadedMs: "1556399443871",
				Name:           "sha256:2066edba3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263",
			},
			{
				ImageSizeBytes: "29905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"master.6d902732dc8c0e19725eaa40c7a860a4c02ef406"},
				TimeCreatedMs:  "1556399434238", // newer by 1ms
				TimeUploadedMs: "1556399443871",
				Name:           "sha256:1366ef2a3485c8f4980ad3d8d96c3ef21de5564d7148146c4978f0d474b67263",
			},
			{
				ImageSizeBytes: "28905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"dev", "latest"},
				TimeCreatedMs:  "1556399434237", // oldest by 1ms
				TimeUploadedMs: "1556399443871",
				Name:           "sha256:5e6a2a225050edcb62cea2a01f4f1e4b2610b6c9e98e8b347f78c49fdc05aff7",
			},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSortedSlice(tt.args.sortBy, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toSortedSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
