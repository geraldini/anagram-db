package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

const AnagramsDbFile = "anagram.json"

func sliceToMap(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	for _, value := range slice {
		set[value] = struct{}{}
	}
	return set
}

func sortWord(word string) string {
	stringsArray := strings.Split(word, "")
	sort.Strings(stringsArray)
	return strings.Join(stringsArray, "")
}

type AnagramDb struct {
	KnownWords map[string][]string
	DbFile     string
}

func (anagramDb *AnagramDb) LoadWords() {
	jsonFile, err := os.Open(anagramDb.DbFile)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &anagramDb.KnownWords)
}

func (anagramDb *AnagramDb) Save() {
	jsonString, err := json.MarshalIndent(anagramDb.KnownWords, "", "  ")
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't save anagrams to file %s.", anagramDb.DbFile))
	}
	log.Printf("Save known words:\n%s", jsonString)
	ioutil.WriteFile(anagramDb.DbFile, jsonString, 0644)
}

func (anagramDb *AnagramDb) GetAnagrams(word string) []string {
	sortedWord := sortWord(word)
	result := []string{}
	if knownWords, ok := anagramDb.KnownWords[sortedWord]; ok {
		log.Printf("There are anagrams for word %s.", word)
		wordsMap := sliceToMap(knownWords)
		if _, ok := wordsMap[word]; ok {
			log.Printf("Word %s was already known.", word)
			for _, value := range anagramDb.KnownWords[sortedWord] {
				if value != word {
					result = append(result, value)
				}
			}
		} else {
			// The word was not known. Add it to the list of known words
			log.Printf("Word %s was not known. Add it to the list of known words.", word)
			result = anagramDb.KnownWords[sortedWord]
			anagramDb.KnownWords[sortedWord] = append(anagramDb.KnownWords[sortedWord], word)
			anagramDb.Save()
		}
	} else {
		log.Printf("There are no anagrams for word %s.", word)
		knownWords := []string{word}
		anagramDb.KnownWords[sortedWord] = knownWords
		anagramDb.Save()
	}
	return result
}

func main() {
	word := flag.String("word", "", "Word to search in the DB.")
	anagramFile := flag.String("file-path", AnagramsDbFile, "File with the known words")
	flag.Parse()
	anagramDb := AnagramDb{DbFile: *anagramFile}
	anagramDb.LoadWords()
	anagrams := anagramDb.GetAnagrams(*word)
	fmt.Printf("These are the known anagrams of %s:\n%s\n\n", *word, anagrams)
}
