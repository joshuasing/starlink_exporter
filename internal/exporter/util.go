// Copyright (c) 2024-2025 Joshua Sing <joshua@joshuasing.dev>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package exporter

import (
	"strconv"
)

// itos converts an int to a string.
func itos[T int | int8 | int16 | int32 | int64](f T) string {
	return strconv.FormatInt(int64(f), 10)
}

// ftos converts a float to a string.
func ftos(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// btof converts the bool to a float64 value of 1/0 (true/false).
func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// parseRingBuffer parses the ringbuffer of history data
func parseRingBuffer(data []float32, current uint64) []float32 {
	bufferSize := uint64(len(data))
	if current <= bufferSize {
		return data[:current]
	}

	start := current % bufferSize
	result := make([]float32, bufferSize)
	copy(result, data[start:])
	copy(result[bufferSize-start:], data[:start])
	return result
}
