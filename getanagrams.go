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

func getKeys(myMap map[string]struct{}, keyToExclude string) []string {
	keys := []string{}
	for key := range myMap {
		if key != keyToExclude {
			keys = append(keys, key)
		}
	}
	return keys
}

func sortWord(word string) string {
	stringsArray := strings.Split(word, "")
	sort.Strings(stringsArray)
	return strings.Join(stringsArray, "")
}

type AnagramDb struct {
	KnownWords map[string]map[string]struct{}
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
	if knownWords, ok := anagramDb.KnownWords[sortedWord]; ok {
		log.Printf("There are anagrams for word %s.", word)
		if _, ok := knownWords[word]; ok {
			log.Printf("Word %s was already known.", word)
		} else {
			// The word was not known. Add it to the list of known words
			log.Printf("Word %s was not known. Add it to the list of known words.", word)
			anagramDb.KnownWords[sortedWord][word] = struct{}{}
			anagramDb.Save()
		}
	} else {
		log.Printf("There are no anagrams for word %s.", word)
		anagramDb.KnownWords[sortedWord] = map[string]struct{}{}
		anagramDb.KnownWords[sortedWord][word] = struct{}{}
		anagramDb.Save()
	}
	return getKeys(anagramDb.KnownWords[sortedWord], word)
}

func main() {
	word := flag.String("word", "", "Word to search in the DB.")
	flag.Parse()
	anagramDb := AnagramDb{DbFile: AnagramsDbFile}
	anagramDb.LoadWords()
	anagrams := anagramDb.GetAnagrams(*word)
	fmt.Printf("These are the known anagrams of %s:\n%s\n\n", *word, anagrams)
}
