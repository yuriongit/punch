package loadtest

import (
	"crypto/rand"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/yuriongit/punch/infra"
	punchTypes "github.com/yuriongit/punch/shared/types"
)

func createTest(config *punchTypes.PunchConfig) (string, error) {
	const MAX_RETRIES uint8 = 5

	for i := 0; i < int(MAX_RETRIES); i++ {
		testId := rand.Text()[0:12]

		existingTestId, err := infra.RdbGetKey(testId)
		if err != nil && err != redis.Nil {
			return "", errors.New("error retrieving test_id for duplication prevention")
		}

		if existingTestId == "" {
			infra.RdbSetKey(testId, *config, 90)
			return testId, nil
		}
	}

	return "", errors.New("failed to generate unique test_id after max retries")
}
