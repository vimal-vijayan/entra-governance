package applications

import (
	"context"
	"sort"
)

func equalStringSets(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	normalizedLeft := append([]string(nil), left...)
	normalizedRight := append([]string(nil), right...)
	sort.Strings(normalizedLeft)
	sort.Strings(normalizedRight)

	for i := range normalizedLeft {
		if normalizedLeft[i] != normalizedRight[i] {
			return false
		}
	}

	return true
}

func findMissing(desired, current []string) []string {
	missing := []string{}
	for _, owner := range desired {
		if !contains(current, owner) {
			missing = append(missing, owner)
		}
	}
	return missing
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func transformHelper(ctx context.Context, appId string, members []string, transform func(ctx context.Context, appId string, members []string) error) error {

	if len(members) > 0 {
		err := transform(ctx, appId, members)
		if err != nil {
			return err
		}
	}
	return nil
}
