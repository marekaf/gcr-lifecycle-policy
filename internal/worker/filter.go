package worker

import (
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
		// `name` contains path in registry

		// split RepositoryPrefix to host and path
		// `gcr.io/repo/` -> [`gcr.io`, `repo/`]
		prefixes := strings.SplitN(repo.RepositoryPrefix, "/", 2)

		// - concatinate repository path and image name
		// - delete `@sha256` suffix from ImageName when specified the image in container_id in k8s manifest
		imageNameInManifest := prefixes[1] + strings.Replace(repo.ImageName, "@sha256", "", 1)

		// check the tags only for the same prefix / image names
		if imageNameInManifest != name {
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

		// if image specified in container_id in k8s manifest
		// compare id instead of the tag name
		if strings.Contains(repo.ImageName, "@sha256") {
			containerID := "sha256:" + repo.Tag
			if containerID == d.Name {
				return true
			}
		}
	}

	return false
}

func digestHasTags(d Digest) bool {
	return len(d.Tag) > 0
}

func protected(pattern string, d Digest) bool {
	for _, tag := range d.Tag {
		if tag == "" {
			continue
		}
		matched, err := regexp.MatchString(pattern, tag)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func filter(c Config, list ListResponse, clusterCat Catalog) FilteredList {

	filtered := FilteredList{}

	keep := c.KeepTags
	retention := daysToTime(time.Now(), c.RetentionDays)

	for _, image := range list.TagsResponses {

		filteredManifests := make(map[string]Digest)

		filteredImage := TagsResponse{
			Name:     image.Name,
			Child:    image.Child,
			Manifest: filteredManifests,
			Tags:     image.Tags,
		}

		sortedDigests := toSortedSlice(c.SortBy, image.Manifest)

		keepCounter := 0

		for _, digest := range sortedDigests {
			if c.ProtectTagRegex != "" {
				if protected(c.ProtectTagRegex, digest) {
					log.Debugf("not deleting digest %+v because match protected pattern (%s)", digest, c.ProtectTagRegex)
					continue
				}
			}

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
			if digestHasTags(digest) && !olderThanRetention(c.SortBy, digest, retention) {
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
