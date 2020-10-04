package worker

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/marekaf/gcr-lifecycle-policy/internal/utils"
)

func duration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder

	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy ", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd", days)

	return b.String()
}

func printBeforeCleanup(list FilteredList) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "REPO", "IMAGE_NAME", "DIGEST", "IMAGE_TAG", "SIZE", "DATE_CREATED", "DATE_UPLOADED"})

	totalSize := 0

	i := 0

	for _, item := range list.TagsResponses {
		for digest, manifest := range item.Manifest {
			repo := extractRepositoryFromImage(item.Name)

			// digest is always prefixed with 'sha256:'
			digestSlug := digest[:27] + "…"

			tagsSlug := strings.Join(manifest.Tag, ",")

			if len(tagsSlug) > 30 {
				tagsSlug = tagsSlug[:27] + "…"
			}

			timecreated, _ := strconv.ParseInt(manifest.TimeCreatedMs, 10, 64)
			timeuploaded, _ := strconv.ParseInt(manifest.TimeUploadedMs, 10, 64)
			ageReadable := time.Unix(timecreated/1000, 0).Format("2006-02-01")
			uploadedReadable := time.Unix(timeuploaded/1000, 0).Format("2006-02-01")

			tmp, _ := strconv.Atoi(manifest.ImageSizeBytes)
			totalSize += tmp

			t.AppendRow([]interface{}{i, repo.RepositoryPrefix, repo.ImageName, digestSlug, tagsSlug, utils.ByteCountSI(manifest.ImageSizeBytes), ageReadable, uploadedReadable})

			i++

			// if i > 100 {
			// 	break
			// }
		}

	}

	t.AppendFooter(table.Row{"", "", "", "", "Total size to save", utils.ByteCountSIInt(totalSize)})
	t.Render()

}

// PrintList prints the report in a pretty table output
func PrintList(list ListResponse) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "REPO", "IMAGE_NAME", "DIGEST", "IMAGE_TAG", "SIZE", "DATE_CREATED", "DATE_UPLOADED"})

	totalSize := 0

	i := 0

	for _, item := range list.TagsResponses {
		for digest, manifest := range item.Manifest {

			repo := extractRepositoryFromImage(item.Name)

			// digest is always prefixed with 'sha256:'
			digestSlug := digest[:27] + "…"

			tagsSlug := strings.Join(manifest.Tag, ",")

			tmp, _ := strconv.Atoi(manifest.ImageSizeBytes)
			totalSize += tmp

			if len(tagsSlug) > 30 {
				tagsSlug = tagsSlug[:27] + "…"
			}

			timecreated, _ := strconv.ParseInt(manifest.TimeCreatedMs, 10, 64)
			timeuploaded, _ := strconv.ParseInt(manifest.TimeUploadedMs, 10, 64)
			ageReadable := time.Unix(timecreated/1000, 0).Format("2006-02-01")
			uploadedReadable := time.Unix(timeuploaded/1000, 0).Format("2006-02-01")

			t.AppendRow([]interface{}{i, repo.RepositoryPrefix, repo.ImageName, digestSlug, tagsSlug, utils.ByteCountSI(manifest.ImageSizeBytes), ageReadable, uploadedReadable})

			i++
		}

	}

	t.AppendFooter(table.Row{"", "", "", "", "Total size", utils.ByteCountSIInt(totalSize)})
	t.Render()
}

// PrintListRepos prints the report in a pretty table output
func PrintListRepos(cat Catalog) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "REPO", "IMAGE_NAME"})

	i := 0

	for _, repo := range cat.Repositories {

		t.AppendRow([]interface{}{i, repo.RepositoryPrefix, repo.ImageName})

		i++
	}

	t.Render()
}

// PrintListCluster prints the report in a pretty table output
func PrintListCluster(cat Catalog) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "REPO", "IMAGE_NAME", "TAG"})

	i := 0

	for _, repo := range cat.Repositories {
		t.AppendRow([]interface{}{i, repo.RepositoryPrefix, repo.ImageName, repo.Tag})

		i++
	}

	t.Render()
}
