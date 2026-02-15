package applications

import "context"

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
