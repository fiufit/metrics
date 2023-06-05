package repositories

import (
	"os"
	"testing"

	testingUtils "github.com/fiufit/metrics/utils/testing"
)

var testSuite testingUtils.TestSuite

const testDbName = "testdb"

func TestMain(m *testing.M) {
	testSuite = testingUtils.NewTestSuite()

	testResult := m.Run()
	testSuite.TearDown()
	os.Exit(testResult)
}
