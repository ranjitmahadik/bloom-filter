package main

import (
	"fmt"
	"log"
	"sync"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/google/uuid"
	"github.com/ranjitmahadik/bloom-filters/core"
)

func generateDataset(size int, train bool) ([]string, map[string]bool) {
	dataset := []string{}
	lookup := map[string]bool{}
	for i := 0; i <= int(size); i++ {
		id := uuid.New()
		dataset = append(dataset, id.String())
		lookup[id.String()] = train
	}
	return dataset, lookup
}

func erroRateForRedis(trainDS, testDS []string, wg *sync.WaitGroup) {
	redisClient := redisbloom.NewClient("localhost:6379", "", nil)
	defer wg.Done()

	redisClient.Reserve("un-sub:users", 0.01, 10_000_000) // at max we can have 10 million users
	_, err := redisClient.BfAddMulti("un-sub:users", trainDS)
	if err != nil {
		log.Fatal("Failed to injest data into bf")
	}
	log.Println("Injested 1 Million users data into BF")

	res, err := redisClient.BfExistsMulti("un-sub:users", testDS)
	if err != nil {
		log.Fatal("Failed to get data from bf")
	}
	falsePositives := 0
	for _, data := range res {
		if data == 1 {
			falsePositives++
		}
	}

	fmt.Println("error rate redis : ", falsePositives)
}

func erroRateForCustomBF(trainDS, testDS []string, wg *sync.WaitGroup) {
	bloomOptions, err := core.NewBloomOptions([]string{"0.01", "10000000"}, false)
	defer wg.Done()
	if err != nil {
		log.Fatalln(err)
		return
	}

	bloomFilter := core.NewBloomFilter(bloomOptions)
	for _, val := range trainDS {
		if _, err := bloomFilter.Add(val); err != nil {
			log.Fatal(err)
			break
		}
	}

	falsePositives := 0
	for _, val := range testDS {
		resp := []byte("-1")
		if resp, err = bloomFilter.Exits(val); err != nil {
			log.Fatal(err)
			break
		}
		if resp[0] == []byte("1")[0] {
			falsePositives++
		}
	}

	fmt.Println("error rate custom : ", falsePositives)
}

func main() {
	dataSetSize := 10_00_000 // 1 Million Users
	trainDs, _ := generateDataset(dataSetSize, true)
	testDs, _ := generateDataset(dataSetSize, false)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go erroRateForRedis(trainDs, testDs, &wg)
	go erroRateForCustomBF(trainDs, testDs, &wg)

	wg.Wait()

}
