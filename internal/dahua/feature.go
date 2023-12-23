package dahua

import (
	"cmp"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func init() {
	FeatureMap = make(map[string]models.DahuaFeature)
	for _, feature := range FeatureList {
		FeatureMap[feature.Key] = feature.DahuaFeature
	}
	slices.SortFunc(FeatureList, func(a Feature, b Feature) int { return cmp.Compare(a.Key, b.Key) })
}

var FeatureList []Feature = []Feature{
	{"camera", "Camera", "", models.DahuaFeatureCamera},
}

type Feature struct {
	Key         string
	Name        string
	Description string
	models.DahuaFeature
}

var FeatureMap map[string]models.DahuaFeature

func FeatureFromStrings(featureStrings []string) models.DahuaFeature {
	var f models.DahuaFeature
	for _, featureString := range featureStrings {
		feature, ok := FeatureMap[featureString]
		if !ok {
			continue
		}
		f = f | feature
	}
	return f
}
