package compute

import (
	"sort"

	"github.com/joyent/triton-go/compute"
)

type imageSort []*compute.Image

func sortImages(images []*compute.Image) []*compute.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
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

type instanceSort []*compute.Instance

func sortInstances(instances []*compute.Instance) []*compute.Instance {
	sortedInstances := instances
	sort.Sort(instanceSort(sortedInstances))
	return sortedInstances
}

func (a instanceSort) Len() int {
	return len(a)
}

func (a instanceSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a instanceSort) Less(i, j int) bool {
	itime := a[i].Created
	jtime := a[j].Created
	return itime.Unix() < jtime.Unix()
}
