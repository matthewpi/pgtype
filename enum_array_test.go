package pgtype_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/matthewpi/pgtype"
	"github.com/matthewpi/pgtype/testutil"
)

func TestEnumArrayTranscode(t *testing.T) {
	setupConn := testutil.MustConnectPgx(t)
	defer testutil.MustCloseContext(t, setupConn)

	if _, err := setupConn.Exec(context.Background(), "drop type if exists color"); err != nil {
		t.Fatal(err)
	}
	if _, err := setupConn.Exec(context.Background(), "create type color as enum ('red', 'green', 'blue')"); err != nil {
		t.Fatal(err)
	}

	testutil.TestSuccessfulTranscode(t, "color[]", []interface{}{
		&pgtype.EnumArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     pgtype.Present,
		},
		&pgtype.EnumArray{
			Elements: []pgtype.GenericText{
				{String: "red", Status: pgtype.Present},
				{Status: pgtype.Null},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		},
		&pgtype.EnumArray{Status: pgtype.Null},
		&pgtype.EnumArray{
			Elements: []pgtype.GenericText{
				{String: "red", Status: pgtype.Present},
				{String: "green", Status: pgtype.Present},
				{String: "blue", Status: pgtype.Present},
				{String: "red", Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: pgtype.Present,
		},
	})
}

func TestEnumArrayArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.EnumArray
	}{
		{
			source: []string{"foo"},
			result: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
		},
		{
			source: (([]string)(nil)),
			result: pgtype.EnumArray{Status: pgtype.Null},
		},
		{
			source: [][]string{{"foo"}, {"bar"}},
			result: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
		},
		{
			source: [][][][]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: pgtype.EnumArray{
				Elements: []pgtype.GenericText{
					{String: "foo", Status: pgtype.Present},
					{String: "bar", Status: pgtype.Present},
					{String: "baz", Status: pgtype.Present},
					{String: "wibble", Status: pgtype.Present},
					{String: "wobble", Status: pgtype.Present},
					{String: "wubble", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: pgtype.Present},
		},
		{
			source: [2][1]string{{"foo"}, {"bar"}},
			result: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
		},
		{
			source: [2][1][1][3]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: pgtype.EnumArray{
				Elements: []pgtype.GenericText{
					{String: "foo", Status: pgtype.Present},
					{String: "bar", Status: pgtype.Present},
					{String: "baz", Status: pgtype.Present},
					{String: "wibble", Status: pgtype.Present},
					{String: "wobble", Status: pgtype.Present},
					{String: "wubble", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: pgtype.Present},
		},
	}

	for i, tt := range successfulTests {
		var r pgtype.EnumArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestEnumArrayArrayAssignTo(t *testing.T) {
	var stringSlice []string
	type _stringSlice []string
	var namedStringSlice _stringSlice
	var stringSliceDim2 [][]string
	var stringSliceDim4 [][][][]string
	var stringArrayDim2 [2][1]string
	var stringArrayDim4 [2][1][1][3]string

	simpleTests := []struct {
		src      pgtype.EnumArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present,
			},
			dst:      &stringSlice,
			expected: []string{"foo"},
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present,
			},
			dst:      &namedStringSlice,
			expected: _stringSlice{"bar"},
		},
		{
			src:      pgtype.EnumArray{Status: pgtype.Null},
			dst:      &stringSlice,
			expected: (([]string)(nil)),
		},
		{
			src:      pgtype.EnumArray{Status: pgtype.Present},
			dst:      &stringSlice,
			expected: []string{},
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
			dst:      &stringSliceDim2,
			expected: [][]string{{"foo"}, {"bar"}},
		},
		{
			src: pgtype.EnumArray{
				Elements: []pgtype.GenericText{
					{String: "foo", Status: pgtype.Present},
					{String: "bar", Status: pgtype.Present},
					{String: "baz", Status: pgtype.Present},
					{String: "wibble", Status: pgtype.Present},
					{String: "wobble", Status: pgtype.Present},
					{String: "wubble", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: pgtype.Present},
			dst:      &stringSliceDim4,
			expected: [][][][]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
			dst:      &stringArrayDim2,
			expected: [2][1]string{{"foo"}, {"bar"}},
		},
		{
			src: pgtype.EnumArray{
				Elements: []pgtype.GenericText{
					{String: "foo", Status: pgtype.Present},
					{String: "bar", Status: pgtype.Present},
					{String: "baz", Status: pgtype.Present},
					{String: "wibble", Status: pgtype.Present},
					{String: "wobble", Status: pgtype.Present},
					{String: "wubble", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: pgtype.Present},
			dst:      &stringArrayDim4,
			expected: [2][1][1][3]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src pgtype.EnumArray
		dst interface{}
	}{
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{Status: pgtype.Null}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     pgtype.Present,
			},
			dst: &stringSlice,
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}, {LowerBound: 1, Length: 2}},
				Status:     pgtype.Present},
			dst: &stringArrayDim2,
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}, {LowerBound: 1, Length: 2}},
				Status:     pgtype.Present},
			dst: &stringSlice,
		},
		{
			src: pgtype.EnumArray{
				Elements:   []pgtype.GenericText{{String: "foo", Status: pgtype.Present}, {String: "bar", Status: pgtype.Present}},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     pgtype.Present},
			dst: &stringArrayDim4,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
