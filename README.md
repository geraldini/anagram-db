# Problem Statement
Two words are anagrams if they have the same letters in a different order, for example, 'stressed' and 'desserts', 'elbow' and 'below'.

Create a program that manages anagrams.
+ When queried for a word, the program should print all the words known by the system that are anagrams of the searched word.
+ The word itself should not be printed.
+ If the word the user searched for isn’t known by the system yet, add that word to the list of known words.

The program should take two arguments:
+ word (mandatory): word to search for anagrams
+ file-path (optional): path of the file where anagrams are stored
