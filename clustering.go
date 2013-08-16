// clustering.go
package main

import (
  "encoding/csv"
  "fmt"
  "clustering-golang/utils"
  "os"
  "strings"
)

// Returns a map where the key is a word and the int is the number
// of times that word appears in the set of documents
// 
// Providing a threshold (1.0 >= x > 0.0) will return only the words
// that appear in all the documents (x*100)% of the time
func termFrequency(fileName string, threshold float32) (m map[string] int, err error) {
  file, err := os.Open(fileName)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  reader := csv.NewReader(file)
  reader.TrailingComma = true

  saveMap := make(map[string] map[string] int)

  recordArray, err := reader.ReadAll()
  if err != nil {
    return nil, err
  }
  for _, record := range recordArray {
    url := record[0]
    
    if _, ok := saveMap[url]; ok {
      continue
    }

    words := utils.LowercaseWords(strings.Fields(record[2]))
  
    for i := range words {
      w, err := utils.RemoveNonAlphaNumeric(words[i])
      if err != nil {
        continue
      } else {
        words[i] = w
      }
    }

    words, err = utils.RemoveStopwords(words)
    if err != nil {
      return nil, err
    }

    saveMap[url] = utils.WordFrequency(words)
  }

  documentFrequencyMap := make(map[string] int)

  for _, wordCountMap := range saveMap {
    for word := range wordCountMap {
      if _, ok := documentFrequencyMap[word]; ok {
        documentFrequencyMap[word]++
      } else {
        documentFrequencyMap[word] = 1
      }
    }
  }

  if threshold != 0.0 {
    for word, value := range documentFrequencyMap {
      if float32(value)/float32(len(saveMap)) < threshold {
        delete(documentFrequencyMap, word)
      }
    }
  }
  return documentFrequencyMap, nil
}

func main() {
  tf, err := termFrequency("pocket.csv", 0.2)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(tf)
  }
}