// The MIT License (MIT)
//
// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package uarand

import (
	"math/rand"
	"time"
)

var (
	// Default is the UARand with default settings.
	Default = New(
		rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
	)
)

// Randomizer represents some entity which could provide us an entropy.
type Randomizer interface {
	Seed(n int64)
	Intn(n int) int
}

// UARand describes the user agent randomizer settings.
type UARand struct {
	Randomizer
}

// GetRandom returns a random user agent from UserAgents slice.
func (u *UARand) GetRandom() string {
	return UserAgents[rand.Intn(len(UserAgents))]
}

// GetRandom returns a random user agent from UserAgents slice.
// This version is driven by Default configuration.
func GetRandom() string {
	return Default.GetRandom()
}

func New(r Randomizer) *UARand {
	return &UARand{r}
}
