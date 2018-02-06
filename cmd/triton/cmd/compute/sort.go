package compute

import (
	"sort"

	"github.com/joyent/triton-go/compute"
)

type imageSort []*compute.Image

func sortImages(images []*compute.Image) []*compute.Image {
	sortedImages := images
	sort.Sort(sort.Reverse(imageSort(sortedImages)))
	return sortedImages
}

func (a imageSort) Len() int {
	return len(a)
}

func (a imageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a imageSort) Less(i, j int) bool {
	itime := a[i].PublishedAt
	jtime := a[j].PublishedAt
	return itime.Unix() < jtime.Unix()
}
