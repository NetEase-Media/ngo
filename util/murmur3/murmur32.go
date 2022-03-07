// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package murmur3

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
)

// MurmurHash copy from guava Murmur3_32HashFunction
type MurmurHash struct {
	seed int
}

// NewMurmurHash constructor of MurmurHash
func NewMurmurHash(seed int) *MurmurHash {
	return &MurmurHash{
		seed: seed,
	}
}

// HashInt32
func (m *MurmurHash) HashInt32(input int32) int {
	k1 := mixK1(input)
	h1 := int32(m.seed)
	h1 = mixH1(h1, k1)
	return int(fmix(h1, 4))
}

// HashInt64
func (m *MurmurHash) HashInt64(input int64) int {
	low := int32(input)
	high := int32(uint(input) >> 32)

	k1 := mixK1(low)
	h1 := int32(m.seed)
	h1 = mixH1(h1, k1)

	k1 = mixK1(high)
	h1 = mixH1(h1, k1)
	return int(fmix(h1, 8))
}

// HashBytes
func (m *MurmurHash) HashBytes(input []byte) int {
	h1 := int32(m.seed)
	iByte := 0
	for ; iByte+4 <= len(input); iByte += 4 {
		k1 := int32(input[iByte]) | int32(input[iByte+1])<<8 | int32(input[iByte+2])<<16 | int32(input[iByte+3])<<24
		k1 = mixK1(k1)
		h1 = mixH1(h1, k1)
	}

	var k1 int32
	switch len(input) - iByte {
	case 3:
		k1 += int32(input[iByte+2]) << 16
		fallthrough
	case 2:
		k1 += int32(input[iByte+1]) << 8
		fallthrough
	case 1:
		k1 += int32(input[iByte])
		k1 = mixK1(k1)
		h1 ^= k1
	}
	return int(fmix(h1, int32(len(input))))
}

func mixK1(k1 int32) int32 {
	k1_64 := int64(k1)
	k1_64 *= c1
	k1 = int32(k1_64)
	k1 = rotateLeft(k1, 15)
	k1_64 = int64(k1)
	k1_64 *= c2
	return int32(k1_64)
}

func mixH1(h1, k1 int32) int32 {
	h1 ^= k1
	h1 = rotateLeft(h1, 13)
	h1_64 := int64(h1)
	h1_64 = h1_64*5 + 0xe6546b64
	return int32(h1_64)
}

func fmix(h1, length int32) int32 {
	uh1 := uint32(h1)
	ulength := uint32(length)
	uh1 ^= ulength
	uh1 ^= uh1 >> 16
	uh1 *= 0x85ebca6b
	uh1 ^= uh1 >> 13
	uh1 *= 0xc2b2ae35
	uh1 ^= uh1 >> 16
	return int32(uh1)
}

func rotateLeft(i, distance int32) int32 {
	ui := uint32(i)
	udistance := uint32(distance)
	a1 := ui << udistance
	b1 := ui >> (32 - udistance)
	return int32(a1 | b1)
}
