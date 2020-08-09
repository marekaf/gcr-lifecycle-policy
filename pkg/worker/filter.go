package worker

import "time"

func filter(list ListResponse) FilteredList {

	filtered := FilteredList{}

	keep := 10
	retention := time.Now().AddDate(-1, 0, 0)

	for _, image := range list.TagsResponses {

		filteredManifests := make(map[string]Digest, 1)

		filteredImage := TagsResponse{
			Name:     image.Name,
			Child:    image.Child,
			Manifest: filteredManifests,
			Tags:     image.Tags,
		}

		sortedManifests := toSortedSlice(image.Manifest)

		keepCounter := 0

		for _, manifest := range sortedManifests {

			if keepCounter < keep {
				keepCounter++
				continue
			}

			// always delete untagged images
			if len(manifest.Tag) > 0 && !olderThanRetention(manifest, retention) {
				continue
			}

			// TODO: we should delete all tags before we delete the image
			filteredManifests[manifest.Name] = Digest(manifest)

		}

		filtered.TagsResponses = append(filtered.TagsResponses, filteredImage)
	}

	return filtered

}
