package worker

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExistsInCluster(t *testing.T) {

	tables := []struct {
		c Catalog
		d Digest
		n string
		z bool
	}{
		{Catalog{
			Repositories: []Repository{
				{
					RepositoryPrefix: "us.gcr.io/mygcp-project/mycorp/",
					ImageName:        "webapp",
					Tag:              "",
				},
				{
					RepositoryPrefix: "gcr.io/mygcp-project/mycorp/",
					ImageName:        "webapp",
					Tag:              "latest",
				},
			},
		}, Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{"master.6d902732dc8c0e19725eaa40c7a860a4c02ef406"},
			TimeCreatedMs:  "1556399434238",
			TimeUploadedMs: "1556399443871",
		}, "mygcp-project/mycorp/webapp", false},
		{Catalog{
			Repositories: []Repository{
				{
					RepositoryPrefix: "us.gcr.io/mygcp-project/mycorp/",
					ImageName:        "webapp",
					Tag:              "",
				},
				{
					RepositoryPrefix: "gcr.io/mygcp-project/mycorp/",
					ImageName:        "webapp",
					Tag:              "latest",
				},
			},
		}, Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{"master.abc", "latest"},
			TimeCreatedMs:  "1556399434238",
			TimeUploadedMs: "1556399443871",
		}, "mygcp-project/mycorp/webapp", true},
		{Catalog{
			Repositories: []Repository{
				{
					RepositoryPrefix: "us.gcr.io/mygcp-project/mycorp/",
					ImageName:        "mysvc",
					Tag:              "master.abc",
				},
				{
					RepositoryPrefix: "gcr.io/mygcp-project/mycorp/",
					ImageName:        "mysvc",
					Tag:              "latest",
				},
			},
		}, Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{"master.abc", "latest"},
			TimeCreatedMs:  "1556399434238",
			TimeUploadedMs: "1556399443871",
		}, "mygcp-project/mycorp/webapp", false},
		{Catalog{
			Repositories: []Repository{
				{
					RepositoryPrefix: "us.gcr.io/mygcp-project/mycorp/",
					ImageName:        "mysvc",
					Tag:              "",
				},
				{
					RepositoryPrefix: "gcr.io/mygcp-project/mycorp/",
					ImageName:        "mysvc",
					Tag:              "latest",
				},
			},
		}, Digest{
			ImageSizeBytes: "29905102",
			LayerID:        "",
			MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
			Tag:            []string{""},
			TimeCreatedMs:  "1556399434238",
			TimeUploadedMs: "1556399443871",
		}, "mygcp-project/mycorp/mysvc", false},
	}

	for _, table := range tables {
		res := existsInCluster(table.c, table.d, table.n)
		diff := cmp.Diff(res, table.z)
		if diff != "" {
			t.Fatalf("existsInCluster mismatch (-want +got):\n%v", diff)
		}
	}
}

func Test_filterCatalog(t *testing.T) {
	type args struct {
		c      Catalog
		filter []string
	}
	tests := []struct {
		name string
		args args
		want Catalog
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterCatalog(tt.args.c, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_protected(t *testing.T) {
	tables := []struct {
		s string
		d Digest
		z bool
	}{
		{"^release-",
			Digest{
				ImageSizeBytes: "29905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"release-20201118-123456"},
				TimeCreatedMs:  "1556399434238",
				TimeUploadedMs: "1556399443871",
			}, true},
		{"^release-",
			Digest{
				ImageSizeBytes: "29905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"test-6ef75b0"},
				TimeCreatedMs:  "1556399434238",
				TimeUploadedMs: "1556399443871",
			}, false},
		{"^v\\d+\\.\\d+\\.\\d+$",
			Digest{
				ImageSizeBytes: "29905102",
				LayerID:        "",
				MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
				Tag:            []string{"v1.2.3"},
				TimeCreatedMs:  "1556399434238",
				TimeUploadedMs: "1556399443871",
			}, true},
	}
	for _, table := range tables {
		res := protected(table.s, table.d)
		diff := cmp.Diff(res, table.z)
		if diff != "" {
			t.Fatalf("protected mismatch (-want +got):\n%v", diff)
		}
	}
}
