package rzerolog

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tmp := "/A/B/C/file.log"
	res := ParseTimeFormat(tmp)
	require.Equal(t, filepath.Join(tmp), res)

	tmp = "/A/B/C/yyyy-MM-dd HH.log"
	res = ParseTimeFormat(tmp)
	require.Equal(t, filepath.Join("/A/B/C/2006-01-02 15.log"), res)
	tmp = "/A/B/C/yyyyyMMMdddHHH.log"
	res = ParseTimeFormat(tmp)
	require.Equal(t, filepath.Join("/A/B/C/2006y01M02d15H.log"), res)
}
