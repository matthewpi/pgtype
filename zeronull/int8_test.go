package zeronull_test

import (
	"testing"

	"github.com/matthewpi/pgtype/testutil"
	"github.com/matthewpi/pgtype/zeronull"
)

func TestInt8Transcode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "int8", []interface{}{
		(zeronull.Int8)(1),
		(zeronull.Int8)(0),
	})
}

func TestInt8ConvertsGoZeroToNull(t *testing.T) {
	testutil.TestGoZeroToNullConversion(t, "int8", (zeronull.Int8)(0))
}

func TestInt8ConvertsNullToGoZero(t *testing.T) {
	testutil.TestNullToGoZeroConversion(t, "int8", (zeronull.Int8)(0))
}
