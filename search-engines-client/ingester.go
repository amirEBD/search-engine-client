package main

import (
	"context"
	"fmt"
	"github.com/expectedsh/go-sonic/sonic"
	"math/rand"
	"sync"
	"time"

	esClient "github.com/olivere/elastic/v7"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

var farsiLetters = []rune("ابپتثجچهخدذرزسشصضطظعغفقهلمنوی")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = farsiLetters[rand.Intn(len(farsiLetters))]
	}
	return string(b)
}

func GenerateFarsiRandomString(n int) string {
	var generatedString string
	for i := 0; i < n; i++ {
		generatedString += RandStringRunes(10) + string(" ")
	}
	return generatedString
}

func NewElasticSearchClient(url, user, pass string) *esClient.Client {
	client, err := esClient.NewClient(
		esClient.SetSniff(false),
		esClient.SetHealthcheck(false),
		esClient.SetURL(url),
		esClient.SetBasicAuth(user, pass))

	if err != nil {
		fmt.Println("client init error:", err)
		return nil
	}

	return client
}

func main() {
	fmt.Println("============ Started =================")
	Init()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()

			es := NewElasticSearchClient("http://localhost:9200", "admin", "testAdmin")
			zinc := NewElasticSearchClient("http://localhost:4080/es", "admin", "testAdmin")
			sonicClient, _ := sonic.NewIngester("localhost", 1491, "password")

			if zinc == nil || es == nil {
				return
			}

			fmt.Println("============ Ingest =================")
			for j := 0; j < 1000000; j++ {
				id := fmt.Sprint(i) + "_" + fmt.Sprint(j)
				indexName := "messages"

				s := GenerateFarsiRandomString(10)
				if rand.Intn(5) == 1 {
					s += "سلام "
				}

				_, err := es.Index().Index(indexName).Id(id).BodyString(fmt.Sprintf(`{"message": "%s"}`, s)).Do(context.Background())
				if err != nil {
					fmt.Println("error-es:", err)
					continue
				}

				_, err = zinc.Index().Index(indexName).Id(id).BodyString(fmt.Sprintf(`{"message": "%s"}`, s)).Do(context.Background())
				if err != nil {
					fmt.Println("error-zinc:", err)
					continue
				}

				err = sonicClient.Push(indexName, "default", id, s, sonic.LangAutoDetect)
				if err != nil {
					fmt.Println("error-sonic:", err)
					continue
				}

			}
		}()
	}
	wg.Wait()

}
