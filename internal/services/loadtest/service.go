// Package loadtest provides load testing services and handlers.
package loadtest

import (
	"crypto/rand"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/yuriongit/punch/internal/infra"
	punchTypes "github.com/yuriongit/punch/shared/types"
)

func createTest(config *punchTypes.PunchConfig) (string, error) {
	const maxRetries uint8 = 5

	for i := 0; i < int(maxRetries); i++ {
		testID := rand.Text()[0:12]

		existingTestID, err := infra.RdbGetKey(testID)
		if err != nil && err != redis.Nil {
			return "", errors.New("error retrieving test_id for duplication prevention")
		}

		if existingTestID == "" {
			if err := infra.RdbSetKey(testID, *config, 90); err != nil {
				return "", errors.New("failed to persist PunchConfig to Redis: " + err.Error())
			}
			return testID, nil
		}
	}

	return "", errors.New("failed to generate unique test_id after max retries")
}
