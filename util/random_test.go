package util

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomDecimalString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "", RandomDecimalString(-1))
	require.Equal(t, "", RandomDecimalString(0))
	n := 10
	for i := 1; i < n; i++ {
		s := RandomDecimalString(i)
		require.Len(t, s, i)
		require.True(t, IsDecimalString(s))
	}

	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = RandomDecimalString(6)
	}

	sort.Strings(a)
	for i := 0; i+1 < n; i++ {
		require.True(t, a[i] < a[i+1])
	}
}

func TestRandomHeximalString(t *testing.T) {
	t.Parallel()

	require.Empty(t, RandomHeximalString(-1))
	require.Empty(t, RandomHeximalString(0))

	n := 10
	for i := 1; i < n; i++ {
		s := RandomHeximalString(i)
		require.Len(t, s, i)
		require.True(t, IsHeximalString(s))
	}

	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = RandomHeximalString(6)
	}

	sort.Strings(a)
	for i := 0; i+1 < n; i++ {
		require.True(t, a[i] < a[i+1])
	}
}

func TestRandomString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Alphabet string
	}{
		{"upper_english", UpperEnglishLetters},
		{"lower_english", LowerEnglishLetters},
		{"decimal_digit", DecimalDigits},
		{"english", EnglishLetters},
		{"alpha_numeric", AlphaNumericCharacters},
		{"unicode_letters", UnicodeLetters},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			for n := 0; n < 20; n++ {
				s := RandomString(n, tc.Alphabet)
				r := []rune(s)
				require.Len(t, r, n)
			}
		})
	}
}

func TestRandomNaturalNumber(t *testing.T) {
	t.Parallel()

	for n := 1; n <= 1000; n++ {
		for i := 0; i <= 100; i++ {
			v := RandomNaturalNumber(n)
			require.True(t, v >= 0)
			require.True(t, v < n)
		}
	}
}
