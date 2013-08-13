package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
  "strings"
  "regexp"
)

func stopwords() (s []string, err error) {
  file, err := os.Open("stopwords.csv")
  if err != nil {
    fmt.Println("Error:", err)
    return nil, err
  }
  defer file.Close()
  reader := csv.NewReader(file)
  record, err := reader.Read()
  if err != nil {
    return nil, err
  } else {
    return record, nil
  }
}

func wordInList(a string, list []string) (b bool, index int) {
  for i, b := range list {
    if b == a {
      return true, i
    }
  }
  return false, -1
}

func lowercaseWords(words []string) []string {
  for i := range words {
    words[i] = strings.ToLower(words[i])
  }
  return words
}

func removeNonAlphaNumeric(s string) (str string, err error) {
  r, err := regexp.Compile("[^\\w]|[-+]?\\d+")
  if err != nil {
    return "", err
  }
  return r.ReplaceAllString(s, ""), nil
}

func documentFrequencyThreshold(fileName string, threshold float32) (m map[string] int, err error) {
  stopwords, err := stopwords()
  if err != nil {
    fmt.Println("Error:", err)
    return nil, err
  }

  file, err := os.Open(fileName)
  if err != nil {
    fmt.Println("Error:", err)
    return nil, err
  }
  defer file.Close()

  reader := csv.NewReader(file)
  reader.TrailingComma = true

  saveMap := make(map[string] map[string] int)

  count := 0

  for {
    record, err := reader.Read()
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Println("Error:", err)
      return nil, err
    }
    count++
    url := record[0]

    if _, ok := saveMap[url]; ok {
      continue
    }

    wordCountMap := make(map[string] int)
    words := lowercaseWords(strings.Fields(record[2]))
  
    for i := range words {
      w, err := removeNonAlphaNumeric(words[i])
      if err != nil {
        continue
      } else {
        words[i] = w
      }
    }

    for _, word := range stopwords { 
      in, i := wordInList(word, words);
      for in {
        words = words[:i+copy(words[i:], words[i+1:])]
        in, i = wordInList(word, words)
      }
    }

    for _, word := range words {
      if strings.EqualFold(word, "") {
        continue
      }
      if _, ok := wordCountMap[word]; ok {
        wordCountMap[word]++
      } else {
        wordCountMap[word] = 1
      }
    }
    saveMap[url] = wordCountMap
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
      if float32(value)/float32(count) < threshold {
        delete(documentFrequencyMap, word)
      }
    }
  }
  return documentFrequencyMap, nil
}

func main() {
  fmt.Println(documentFrequencyThreshold("pocket.csv", 0.5))
}