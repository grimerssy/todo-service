package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFU(t *testing.T) {
	tests := []struct {
		cfg      ConfigLFU
		testCase func(c *LFU) []interface{}
		want     []interface{}
	}{
		{
			cfg: ConfigLFU{
				capacity: 0,
			},
			testCase: func(c *LFU) []interface{} {
				var results []interface{}

				c.SetValue(0, 0)
				results = append(results, c.GetValue(0))

				return results
			},
			want: []interface{}{
				-1,
			},
		},
		{
			cfg: ConfigLFU{
				capacity: 3,
			},
			testCase: func(c *LFU) []interface{} {
				var results []interface{}

				c.SetValue(2, 2)
				c.SetValue(1, 1)
				results = append(results, c.GetValue(2))
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(2))
				c.SetValue(3, 3)
				c.SetValue(4, 4)
				results = append(results, c.GetValue(3))
				results = append(results, c.GetValue(2))
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(4))

				return results
			},
			want: []interface{}{
				2, 1, 2, -1, 2, 1, 4,
			},
		},
		{
			cfg: ConfigLFU{
				capacity: 2,
			},
			testCase: func(c *LFU) []interface{} {
				var results []interface{}

				c.SetValue(1, 1)
				c.SetValue(2, 2)
				results = append(results, c.GetValue(1))
				c.SetValue(3, 3)
				results = append(results, c.GetValue(2))
				results = append(results, c.GetValue(3))
				c.SetValue(4, 4)
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(3))
				results = append(results, c.GetValue(4))

				return results
			},
			want: []interface{}{
				1, -1, 3, -1, 3, 4,
			},
		},
	}
	for _, tt := range tests {
		c := NewLFU(tt.cfg)
		got := tt.testCase(c)
		assert.Equal(t, tt.want, got)
	}
}
