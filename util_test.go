package xmlpath

import "testing"

func Test_TestInterference(t *testing.T) {
	var callback Decoder = func(decodeInto func(i interface{})) {}
	path1 := NewPathConfig(callback, "alfa", "beta", "gamma", "delta")
	path2 := NewPathConfig(callback, "alfa", "beta", "theta")
	path3 := NewPathConfig(callback, "alfa", "beta", "theta", "sigma")
	tt := []struct {
		arg          []PathConfig
		resultTester func(e error) bool
	}{
		{
			[]PathConfig{path1, path2},
			func(e error) bool { return e == nil },
		},
		{
			[]PathConfig{path1, path2, path3},
			func(e error) bool { return e != nil },
		},
		{
			[]PathConfig{path2, path3},
			func(e error) bool { return e != nil },
		},
	}
	for i, tc := range tt {
		if !tc.resultTester(testInterference(tc.arg)) {
			t.Errorf("Failed test #%v", i)
		}
	}
}
