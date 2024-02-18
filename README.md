# Bloom Filter in Go

This repository contains a simple implementation of a Bloom filter data structure in Go.

## Introduction

A Bloom filter is a space-efficient probabilistic data structure used to test whether an element is a member of a set. It allows for quick membership queries with a controlled probability of false positives. While false positive matches are possible, false negatives are not. 

## Setup

To get started, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/ranjitmahadik/bloom-filter.git
2. **Navigate to the cloned directory:**
   ```bash
   cd bloom-filter
2. **Now run the following command:**
   ```bash
   make run

## **Usage**
Here's how you can use the Bloom filter.
1. Initialize the Bloom filter options with the desired error rate and capacity:
```go
  bloomOptions, err := core.NewBloomOptions([]string{"0.01", "10000000"}, false)
```
Replace **"0.01"** with the desired error rate and **"10000000"** with the desired capacity.

2. Now initialize the bloom filter with bloom filter options.
```go
  bloomOptions, err := core.NewBloomOptions([]string{"0.01", "10000000"}, false)
  bloomFilter := core.NewBloomFilter(bloomOptions)
```
3. How to add item to Bloom Filter?
```go
  bloomOptions, err := core.NewBloomOptions([]string{"0.01", "10000000"}, false)
  bloomFilter := core.NewBloomFilter(bloomOptions)
  resp, err := bloomFilter.Add("A");
```
4. How to check if value exists in bloom filter?
```go
  bloomOptions, err := core.NewBloomOptions([]string{"0.01", "10000000"}, false)
  bloomFilter := core.NewBloomFilter(bloomOptions)
  resp, err := bloomFilter.Add("A");
  resp, err = bloomFilter.Exits("A");
  if resp[0] == []byte("1")[0] {
    fmt.Println("value present in bloom filter.")
  }
```
