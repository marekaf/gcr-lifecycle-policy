package worker

import log "github.com/sirupsen/logrus"

func filterCatalog(c Catalog, filter []string) Catalog {

	// by default do not filter anything
	if len(filter) == 0 {
		return c
	}

	filtered := Catalog{}

	for _, repo := range c.Repositories {

		for _, filterRepo := range filter {
			if (repo.RepositoryPrefix + repo.ImageName) == filterRepo {
				log.Debugf("repo %s matched the filter", filterRepo)
				filtered.Repositories = append(filtered.Repositories, repo)
				break
			}

			log.Debugf("repo %s did not match the filter %s", repo.RepositoryPrefix+repo.ImageName, filterRepo)
		}
	}

	return filtered
}

func existsInCluster(c Catalog, d Digest, name string) bool {

	for _, repo := range c.Repositories {
		// check the tags only for the same prefix / image names
		if repo.ImageName != name {
			continue
		}
		for _, tag := range d.Tag {

			// don't compare empty tags, skip them
			if tag == "" {
				continue
			}

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
				log.Debugf("not deleting digest %+v because keep-tags (%d/%d)", digest, keepCounter, keep)
				continue
			}

			if existsInCluster(clusterCat, digest, image.Name) {
				log.Debugf("not deleting digest %+v because existsInCluster", digest)
				continue
			}

			// always delete untagged images
			if digestHasTags(digest) && !olderThanRetention(digest, retention) {
				log.Debugf("not deleting digest %+v because digestHasTags && !olderThanRetention (%s)", digest, retention)
				continue
			}

			log.Debugf("adding digest %+v to cleanupList", digest)
			// TODO: we should delete all tags before we delete the image
			filteredManifests[digest.Name] = Digest(digest)

		}

		filtered.TagsResponses = append(filtered.TagsResponses, filteredImage)
	}

	return filtered

}
