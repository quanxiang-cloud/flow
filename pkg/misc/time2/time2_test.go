package time2

import (
	"testing"

	"github.com/bmizerany/assert"
	. "github.com/quanxiang-cloud/flow/pkg/misc/test"
	"github.com/stretchr/testify/suite"
)

func (suite *timeTestSuite) TestTime() {
	nowUnix := NowUnix()
	assert.NotEqual(suite.T(), 0, nowUnix)

	now := Now()
	assert.NotEqual(suite.T(), "", nowUnix)

	ntnu, err := ISO8601ToUnix(now)
	assert.Equal(suite.T(), nil, err)
	assert.NotEqual(suite.T(), 0, ntnu)

	tnun := UnixToISO8601(nowUnix)
	assert.NotEqual(suite.T(), "", tnun)

	mill := NowUnixMill()
	assert.NotEqual(suite.T(), 0, mill)
}

type timeTestSuite struct {
	Suite
}

func TestTimeTestSuite(t *testing.T) {
	suite.Run(t, new(timeTestSuite))
}
