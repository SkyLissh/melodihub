package utils

import "strings"

// Function that gives a numeric rank based on two string and how much they match. The rank is between 0 and 100, where 100 means a perfect match and 0 means no match at all.
func RankString(str1, str2 string) int {
	// Convert both strings to lowercase for case-insensitive comparison
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)

	// Split the strings into words
	words1 := strings.Fields(str1)
	words2 := strings.Fields(str2)

	// Create a map to count the occurrences of each word in str1
	wordCount := make(map[string]int)
	for _, word := range words1 {
		wordCount[word]++
	}

	// Count the number of matching words in str2
	matchCount := 0
	for _, word := range words2 {
		if count, exists := wordCount[word]; exists && count > 0 {
			matchCount++
			wordCount[word]-- // Decrease the count to avoid counting duplicates
		}
	}

	// Calculate the rank as a percentage of matching words
	totalWords := len(words1) + len(words2)
	if totalWords == 0 {
		return 100 // Both strings are empty, consider it a perfect match
	}

	rank := (2 * matchCount * 100) / totalWords // Multiply by 2 to account for both strings
	return rank
}
