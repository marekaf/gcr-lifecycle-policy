package worker

import (
	"fmt"
	"time"
)

func filterCatalog(c Catalog, filter []string) Catalog {

	// by default do not filter anything
	if len(filter) == 0 {
		return c
	}

	filtered := Catalog{}

	for _, repo := range c.Repositories {

		for _, filterRepo := range filter {
			if (repo.RepositoryPrefix + repo.ImageName) == filterRepo {
				filtered.Repositories = append(filtered.Repositories, repo)
				break
			}
		}
	}

	return filtered
}

func daysToTime(days int) time.Time {
	return time.Now().AddDate(0, -days, 0)
}

func existsInCluster(c Catalog, d Digest) bool {

	for _, repo := range c.Repositories {
		for _, tag := range d.Tag {
			if repo.Tag == tag {
				return true
			}
		}
	}

	return false
}

func digestHasTags(d Digest) bool {
	return len(d.Tag) > 0
}

func filter(c Config, list ListResponse, clusterCat Catalog) FilteredList {

	filtered := FilteredList{}

	keep := c.KeepTags
	retention := daysToTime(c.RetentionDays)

	for _, image := range list.TagsResponses {

		filteredManifests := make(map[string]Digest, 0)

		filteredImage := TagsResponse{
			Name:     image.Name,
			Child:    image.Child,
			Manifest: filteredManifests,
			Tags:     image.Tags,
		}

		sortedDigests := toSortedSlice(image.Manifest)

		keepCounter := 0

		for _, digest := range sortedDigests {

			if keepCounter < keep {
				keepCounter++
				//fmt.Printf("not deleting digest %+v because keep-tags (%d/%d)", digest, keepCounter, keep)
				continue
			}

			if existsInCluster(clusterCat, digest) {
				//fmt.Printf("not deleting digest %+v because existsInCluster", digest)
				continue
			}

			// always delete untagged images
			if digestHasTags(digest) && !olderThanRetention(digest, retention) {
				if digestHasTags(digest) {
					//fmt.Printf("not deleting digest %+v because digestHasTags", digest)
				} else {
					//fmt.Printf("not deleting digest %+v because !olderThanRetention (%s)", digest, retention)
				}
				continue
			}

			fmt.Printf("deleting digest %+v", digest)
			// TODO: we should delete all tags before we delete the image
			filteredManifests[digest.Name] = Digest(digest)

		}

		filtered.TagsResponses = append(filtered.TagsResponses, filteredImage)
	}

	return filtered

}
