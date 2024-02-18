package core

import (
	"errors"
	"fmt"
	"hash"
	"math"
	"math/rand"
	"strconv"

	"github.com/spaolacci/murmur3"
)

const (
	defaultErrorRate float64 = 0.01
	defaultCapacity  uint64  = 1024
)

var (
	ln2      float64 = math.Log(2)
	ln2Power float64 = ln2 * ln2
)

var (
	invalidErrorType = errors.New("ERR: only float64 error type allowed for error rate.")
	invalidErrorRate = errors.New("ERR: error rate should be greater than 0 and less than 1.")

	invalidCapacityType = errors.New("ERR: only uint64 capacity type allowed for capacity.")
	invalidCapacity     = errors.New("ERR: capacity should be greater than 1.")

	emptyValueError   = errors.New("ERR: empty value can't be added to bloom filter")
	unableToHashError = errors.New("ERR: unable hash value")
)

type BloomOptions struct {
	errorRate float64       // desired error rate
	capacity  uint64        // number of entries to be added in bloom filter
	bits      uint64        // total number of bits reserved in bloom filter
	hashFns   []hash.Hash64 // list of hash functions
	bpe       float64       // bits per element

	// maintains hashed indexes for all the hash functions.
	indexes []uint64
}

type Bloom struct {
	options *BloomOptions // options for bf
	filter  []byte        // bit representation over array.
}

func NewBloomOptions(args []string, useDefaults bool) (*BloomOptions, error) {
	if useDefaults {
		return &BloomOptions{
				errorRate: defaultErrorRate,
				capacity:  defaultCapacity,
			},
			nil
	}

	errorRate, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, invalidErrorType
	}

	if errorRate <= 0 || errorRate >= 1.0 {
		return nil, invalidErrorRate
	}

	capacity, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return nil, invalidCapacityType
	}

	if capacity < 1 {
		return nil, invalidCapacity
	}

	return &BloomOptions{
		errorRate: errorRate,
		capacity:  capacity,
	}, nil
}

func NewBloomFilter(options *BloomOptions) *Bloom {
	/**
	calculate bits per element
	bpe = -log(errorRate)/log(2)^2
	*/
	options.bpe = (-1 * math.Log(options.errorRate)) / ln2Power

	/**
	calculate the required hash functions (k)
	k = ceil(ln(2) * bpe)
	*/
	k := int(math.Ceil(ln2 * options.bpe))
	options.hashFns = make([]hash.Hash64, k)

	for i := 0; i < k; i++ {
		options.hashFns[i] = murmur3.New64WithSeed(rand.Uint32())
	}

	options.indexes = make([]uint64, k)

	/**
	calculate number of required bits
	bits := k * entries /ln(2)
	bytes = bits * 8
	*/
	bits := uint64(math.Ceil((float64(k) * float64(options.capacity)) / ln2))
	bytes := uint64(math.Ceil(float64(bits) / 8))
	options.bits = bits
	bitset := make([]byte, bytes)

	return &Bloom{
		options: options,
		filter:  bitset,
	}
}

func (bf *Bloom) Info(key string) string {
	info := ""
	if key != "" {
		info = "key: " + key + ", "
	}
	info += fmt.Sprintf("error rate : %f\n", bf.options.errorRate)
	info += fmt.Sprintf("capacity : %d\n", bf.options.capacity)
	info += fmt.Sprintf("total bits reserved : %d\n", bf.options.bits)
	info += fmt.Sprintf("bits per element : %f\n", bf.options.bpe)
	info += fmt.Sprintf("hash functions : %d", len(bf.options.hashFns))
	return info
}

func (bf *Bloom) Add(value string) ([]byte, error) {
	if value == "" {
		return []byte("-1"), emptyValueError
	}

	if err := bf.options.updateIndexes(value); err != nil {
		return []byte("-1"), unableToHashError
	}

	var allBitsSetAlready bool = true

	for _, v := range bf.options.indexes {
		if !isBitSet(bf.filter, v) {
			allBitsSetAlready = false
			setBit(bf.filter, v)
		}
	}

	if allBitsSetAlready {
		return []byte("0"), nil
	}

	return []byte("1"), nil

}

func (bf *Bloom) Exits(value string) ([]byte, error) {
	if value == "" {
		return []byte("-1"), emptyValueError
	}

	if err := bf.options.updateIndexes(value); err != nil {
		return []byte("-1"), unableToHashError
	}

	for _, v := range bf.options.indexes {
		if !isBitSet(bf.filter, v) {
			return []byte("-1"), nil
		}
	}

	return []byte("1"), nil
}

func (opts *BloomOptions) updateIndexes(value string) error {
	for i := 0; i < len(opts.hashFns); i++ {
		currentHashFunction := opts.hashFns[i]
		currentHashFunction.Reset()

		if _, err := currentHashFunction.Write([]byte(value)); err != nil {
			return err
		}
		opts.indexes[i] = currentHashFunction.Sum64() % opts.bits
	}

	return nil
}
