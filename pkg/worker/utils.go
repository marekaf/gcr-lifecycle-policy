package worker

import "strings"

func extractRepositoryFromImage(input string) Repository {

	prefix := ""
	image := ""
	tag := ""

	// +1 to keep the trailing slash
	whereToSplit := strings.LastIndex(input, "/") + 1
	prefix = input[:whereToSplit]
	imagetag := input[whereToSplit:]

	// no +1 because we don't want the ':'
	wts := strings.LastIndex(imagetag, ":")
	if wts == -1 {
		image = imagetag
	} else {
		image = imagetag[:wts]
		tag = imagetag[wts+1:]
	}

	repo := Repository{
		RepositoryPrefix: prefix,
		ImageName:        image,
		Tag:              tag,
	}

	return repo
}
