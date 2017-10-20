package tests

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// Алфавит
type Abc struct {
	alphabet []rune
	bits     float64
}

// Идентификатор
type Shortid struct {
	abc   Abc
	rnd   *rand.Rand
	epoch time.Time  // Начало времен
	ms    uint       // Время генерации последнего значения
	count uint       // Кол-во идентификаторов в рамках одой и той же мс.
	mx    sync.Mutex // Блокировка для конкурентного доступа
}

// Создание нового генератора
func NewShortid(alphabet string) *Shortid {
	return &Shortid{
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
		abc:   NewAbc(alphabet),
		epoch: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
		ms:    0,
		count: 0,
	}
}

func (sid *Shortid) Generate() string {
	for true {
		idrunes := sid.GenerateRandom()
		valid := false
		for c := range idrunes {
			if (c >= 'a') || (c <= 'z') {
				valid = true
				break
			}
		}
		if valid {
			rnd := sid.rnd.Int()
			flg := 1
			for i, c := range idrunes {
				if (c >= 'a') && (c <= 'z') {
					if rnd&flg == 0 {
						idrunes[i] = c + 'A' - 'a'
					}
					flg <<= 1
					if flg == 0 {
						flg = 1
					}
				}
			}
			return string(idrunes)
		}
	}
	panic("Unreacheable code")
}

// Generate generates a new short Id.
func (sid *Shortid) GenerateRandom() []rune {
	sid.mx.Lock()
	defer sid.mx.Unlock()

	ms, count := sid.getMsAndCounter(sid.epoch)
	idrunes := sid.abc.Encode(sid.rnd, ms, 40)
	if count > 0 {
		idrunes = append(idrunes, sid.abc.Encode(sid.rnd, count, 0)...)
	}
	return idrunes
}

func (sid *Shortid) getMsAndCounter(epoch time.Time) (uint, uint) {
	ms := uint(time.Now().Sub(epoch).Nanoseconds() / 1000000)
	if ms <= sid.ms {
		sid.count++
	} else {
		sid.count = 0
		sid.ms = ms
	}
	return sid.ms, sid.count
}

// Abc returns the instance of alphabet used for representing the Ids.
func (sid *Shortid) Abc() Abc {
	return sid.abc
}

// Epoch returns the value of epoch used as the beginning of millisecond counting (normally
// 2016-01-01 00:00:00 local time)
func (sid *Shortid) Epoch() time.Time {
	return sid.epoch
}

// NewAbc constructs a new instance of shuffled alphabet to be used for Id representation.
func NewAbc(alphabet string) Abc {
	abc := Abc{alphabet: nil, bits: math.Log2(float64(len(alphabet)))}
	abc.shuffle(alphabet, 0)
	return abc
}

func (abc *Abc) shuffle(alphabet string, seed uint64) {
	source := []rune(alphabet)
	for len(source) > 1 {
		seed = (seed*9301 + 49297) % 233280
		i := int(seed * uint64(len(source)) / 233280)

		abc.alphabet = append(abc.alphabet, source[i])
		source = append(source[:i], source[i+1:]...)
	}
	abc.alphabet = append(abc.alphabet, source[0])
}

// Encode encodes a given value into a slice of runes of length nsymbols. In case nsymbols==0, the
// length of the result is automatically computed from data. Even if fewer symbols is required to
// encode the data than nsymbols, all positions are used encoding 0 where required to guarantee
// uniqueness in case further data is added to the sequence. The value of digits [4,6] represents
// represents n in 2^n, which defines how much randomness flows into the algorithm: 4 -- every value
// can be represented by 4 symbols in the alphabet (permitting at most 16 values), 5 -- every value
// can be represented by 2 symbols in the alphabet (permitting at most 32 values), 6 -- every value
// is represented by exactly 1 symbol with no randomness (permitting 64 values).
func (abc *Abc) Encode(rnd *rand.Rand, val, bits uint) []rune {
	nsymbols := uint(0)
	randBits := uint(2)
	if bits > 0 {
		nsymbols = uint(math.Ceil(float64(bits) / (abc.bits - float64(randBits))))
	} else if val > 0 {
		nsymbols = uint(math.Ceil(math.Log2(float64(val+1)) / (abc.bits - float64(randBits))))
	}
	if nsymbols == 0 {
		return []rune{}
	}
	// no random component if digits == 6
	res := make([]rune, int(nsymbols))
	data := val
	for i := range res {
		index := data % uint(len(abc.alphabet)>>randBits)
		index = (index << randBits) | uint(rnd.Int31n(1<<randBits))
		res[i] = abc.alphabet[index]
		data /= uint(len(abc.alphabet) >> randBits)
	}
	return res
}
