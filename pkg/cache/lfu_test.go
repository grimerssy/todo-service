package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFU(t *testing.T) {
	const key = TodoKey

	tests := []struct {
		cfg      ConfigLFU
		testCase func(c *LFU) []any
		want     []any
	}{
		{
			cfg: ConfigLFU{
				Capacities: map[cfgKey]int{
					key: 3,
				},
				CleanupSizes: map[cfgKey]int{
					key: 1,
				},
			},
			testCase: func(c *LFU) []any {
				var results []any

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
			want: []any{
				2, 1, 2, nil, 2, 1, 4,
			},
		},
		{
			cfg: ConfigLFU{
				Capacities: map[cfgKey]int{
					key: 2,
				},
				CleanupSizes: map[cfgKey]int{
					key: 1,
				},
			},
			testCase: func(c *LFU) []any {
				var results []any

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
			want: []any{
				1, nil, 3, nil, 3, 4,
			},
		},
		{
			cfg: ConfigLFU{
				Capacities: map[cfgKey]int{
					key: 1,
				},
				CleanupSizes: map[cfgKey]int{
					key: 1,
				},
			},
			testCase: func(c *LFU) []any {
				var results []any

				c.SetValue(1, 1)
				c.SetValue(2, 2)
				c.SetValue(3, 3)
				results = append(results, c.GetValue(3))
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(2))
				c.RemoveValue(3)
				results = append(results, c.GetValue(3))

				return results
			},
			want: []any{
				3, nil, nil, nil,
			},
		},
		{
			cfg: ConfigLFU{
				Capacities: map[cfgKey]int{
					key: 3,
				},
				CleanupSizes: map[cfgKey]int{
					key: 1,
				},
			},
			testCase: func(c *LFU) []any {
				var results []any

				c.SetValue(1, 1)
				c.SetValue(2, 2)
				c.SetValue(3, 3)
				c.SetValue(4, 4)
				c.SetValue(5, 5)
				results = append(results, c.GetValue(3))
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(4))
				c.RemoveValue(4)
				results = append(results, c.GetValue(4))
				results = append(results, c.GetValue(3))

				return results
			},
			want: []any{
				3, nil, 4, nil, 3,
			},
		},
		{
			cfg: ConfigLFU{
				Capacities: map[cfgKey]int{
					key: 2,
				},
				CleanupSizes: map[cfgKey]int{
					key: 2,
				},
			},
			testCase: func(c *LFU) []any {
				var results []any

				c.SetValue(1, 1)
				c.SetValue(2, 2)
				results = append(results, c.GetValue(2))
				results = append(results, c.GetValue(1))
				c.SetValue(3, 3)
				results = append(results, c.GetValue(2))
				results = append(results, c.GetValue(1))
				results = append(results, c.GetValue(3))
				c.SetValue(4, 4)
				results = append(results, c.GetValue(4))
				results = append(results, c.GetValue(3))
				c.RemoveValue(3)
				c.SetValue(5, 5)
				results = append(results, c.GetValue(5))

				return results
			},
			want: []any{
				2, 1, nil, nil, 3, 4, 3, 5,
			},
		},
	}
	for _, tt := range tests {
		c := NewLFU(tt.cfg, key)
		got := tt.testCase(c)
		assert.Equal(t, tt.want, got)
	}
}
