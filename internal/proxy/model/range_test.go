/**
* Copyright 2018 Comcast Cable Communications Management, LLC
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package model

import (
	"net/http"
	"testing"
)

func TestRanges_CalculateDelta_FullCacheMiss(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 15, End: 20}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 15 ||
		res[0].End != 20 {
		t.Errorf("expected start %d end %d, got start %d end %d", 15, 20, res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_FullCacheMiss2(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 1, End: 4}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 1 ||
		res[0].End != 4 {
		t.Errorf("expected start %d end %d, got start %d end %d", 1, 4, res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_FullCacheMiss3(t *testing.T) {
	ranges := make(Ranges, 2)
	ranges[0] = Range{Start: 5, End: 10}
	ranges[1] = Range{Start: 18, End: 20}
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	// Query
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 12, End: 16}
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 12 ||
		res[0].End != 16 {
		t.Errorf("expected start %d end %d, got start %d end %d", 12, 16, res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 0, End: 0}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of out of bounds, got start %d end %d", res[0].Start, res[0].End)
	}
	ranges2[0] = Range{Start: 100, End: 100}
	res = ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of out of bounds, got start %d end %d", res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_EmptyContentRange(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 5, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}
}

func TestRanges_CalculateDelta_InvalidContentRange(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 5, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}
}

func TestRanges_CalculateDelta_InvalidContentRange2(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/blah")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 5, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}
}

func TestRanges_CalculateDelta_OutOfBounds(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: -1, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}
	ranges2[0] = Range{Start: 1, End: 100}
	res = ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}

	ranges2[0] = Range{Start: -10, End: 65}
	res = ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of empty content range, but got a valid value")
	}
}

func TestRanges_CalculateDelta_EmptyContentLength(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 5, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of out of bounds, got start %d end %d", res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_InvalidContentLength(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62thisiswrong")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)

	// Queries
	ranges2 := make(Ranges, 1)
	ranges2[0] = Range{Start: 5, End: 10}
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected nil because of out of bounds, got start %d end %d", res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_PartialCacheMiss(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 8, End: 20}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 10 ||
		res[0].End != 20 {
		t.Errorf("expected start %d end %d, got start %d end %d", 10, 20, res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_PartialCacheMiss2(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 2, End: 8}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 2 ||
		res[0].End != 5 {
		t.Errorf("expected start %d end %d, got start %d end %d", 2, 5, res[0].Start, res[0].End)
	}
}

func TestRanges_CalculateDelta_PartialCacheMiss3(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 2, End: 20}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 2 ||
		res[0].End != 4 {
		t.Errorf("expected start %d end %d, got start %d end %d", 2, 4, res[0].Start, res[0].End)
	}
	if res[1].Start != 11 ||
		res[1].End != 20 {
		t.Errorf("expected start %d end %d, got start %d end %d", 11, 20, res[1].Start, res[1].End)
	}
}

func TestRanges_CalculateDelta_PartialCacheMiss4(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 2)
	ranges[0] = byteRange
	byteRange2 := Range{Start: 15, End: 20}
	ranges[1] = byteRange2

	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.Header.Add("Content-Range", "bytes 5-10/62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 2, End: 25}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res[0].Start != 2 ||
		res[0].End != 4 {
		t.Errorf("expected start %d end %d, got start %d end %d", 2, 4, res[0].Start, res[0].End)
	}
	if res[1].Start != 11 ||
		res[1].End != 14 {
		t.Errorf("expected start %d end %d, got start %d end %d", 11, 20, res[1].Start, res[1].End)
	}
	if res[2].Start != 21 ||
		res[2].End != 25 {
		t.Errorf("expected start %d end %d, got start %d end %d", 21, 25, res[2].Start, res[2].End)
	}
}

func TestRanges_CalculateDelta_CacheHit(t *testing.T) {
	byteRange := Range{Start: 5, End: 10}
	ranges := make(Ranges, 1)
	ranges[0] = byteRange
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Length", "62")
	resp.StatusCode = 200
	d := DocumentFromHTTPResponse(resp, []byte("This is a test file, to see how the byte range requests work.\n"), nil)
	byteRange = Range{Start: 6, End: 9}
	ranges2 := make(Ranges, 1)
	ranges2[0] = byteRange
	res := ranges.CalculateDelta(d, ranges2)
	if res != nil {
		t.Errorf("expected cache hit but got cache miss")
	}
}

func TestGetByteRanges_EmptyString(t *testing.T) {
	r := GetByteRanges("")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
}

func TestGetByteRanges_InvalidRange(t *testing.T) {
	r := GetByteRanges("bytes=abc-def")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("bytes0-100")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("0-100")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("100")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("-")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("bytes=20-30-40-50")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
	r = GetByteRanges("bytes=20-blah")
	if r != nil {
		t.Errorf("expected empty byte range")
	}
}

func TestGetByteRanges_SingleByteRange(t *testing.T) {
	byteRange := "bytes=0-50"
	res := GetByteRanges(byteRange)
	if res == nil {
		t.Errorf("expected a non empty byte range, but got an empty range")
	}
	if res[0].Start != 0 || res[0].End != 50 {
		t.Errorf("expected start %d end %d, got start %d end %d", 0, 50, res[0].Start, res[0].End)
	}
}

func TestGetByteRanges_MultiByteRange(t *testing.T) {
	byteRange := "bytes=0-50, 100-150"
	res := GetByteRanges(byteRange)
	if res == nil {
		t.Errorf("expected a non empty byte range, but got an empty range")
	}
	if res[0].Start != 0 || res[0].End != 50 {
		t.Errorf("expected start %d end %d, got start %d end %d", 0, 50, res[0].Start, res[0].End)
	}
	if res[1].Start != 100 || res[1].End != 150 {
		t.Errorf("expected start %d end %d, got start %d end %d", 100, 150, res[1].Start, res[1].End)
	}
}