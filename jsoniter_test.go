package jsoniter_test

import (
	"encoding/json"
	"errors"
	"github.com/simon-engledew/jsoniter"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	matcher := jsoniter.Matcher("some", 0, "nested", "structure")

	var found any

	fn := func(path []json.Token) error {
		if matcher(path) {
			return d.Decode(&found)
		}
		return nil
	}

	require.NoError(t, jsoniter.Iterate(d, fn))
	require.Equal(t, map[string]any{"a": 1.0}, found)
}

func TestInvalid(t *testing.T) {
	doc := `{
	  "some": [}`

	d := json.NewDecoder(strings.NewReader(doc))

	fn := func(path []json.Token) error {
		return nil
	}

	require.ErrorContains(t, jsoniter.Iterate(d, fn), `invalid character '}' looking for beginning of value`)
}

func TestEOF(t *testing.T) {
	doc := `{
	  "some": [{`

	d := json.NewDecoder(strings.NewReader(doc))

	fn := func(path []json.Token) error {
		return nil
	}

	require.ErrorIs(t, jsoniter.Iterate(d, fn), io.EOF)
}

func TestIterate(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }, {
		"nested": {
		  "structure": {
			"b": 2
		  }
		}
      }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	matcher := jsoniter.Matcher("some", jsoniter.Wildcard, "nested", "structure")

	var hits int

	fn := func(path []json.Token) error {
		if matcher(path) {
			hits += 1
		}
		return nil
	}

	require.NoError(t, jsoniter.Iterate(d, fn))
	require.Equal(t, 2, hits)
	require.Equal(t, d.InputOffset(), int64(len(doc)))
}

func TestStop(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	stopErr := errors.New("stop")

	fn := func(path []json.Token) error {
		return stopErr
	}

	require.ErrorIs(t, jsoniter.Iterate(d, fn), stopErr)
	require.Less(t, d.InputOffset(), int64(len(doc)))
}
