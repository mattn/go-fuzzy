package fuzzy

import (
	"fmt"
	"unicode"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Match(pattern, str string) (bool, int, string) {
	// Score consts
	var adjacency_bonus = 5             // bonus for adjacent matches
	var separator_bonus = 10            // bonus if match occurs after a separator
	var camel_bonus = 10                // bonus if match is uppercase and prev is lower
	var leading_letter_penalty = -3     // penalty applied for every letter in str before the first match
	var max_leading_letter_penalty = -9 // maximum penalty for leading letters
	var unmatched_letter_penalty = -1   // penalty for every letter that doesn't matter

	// Loop variables
	var score = 0
	var patternIdx = 0
	var patternLength = len([]rune(pattern))
	var strIdx = 0
	var strLength = len([]rune(str))
	var prevMatched = false
	var prevLower = false
	var prevSeparator = true // true so if first letter match gets separator bonus

	// Use "best" matched letter if multiple string letters match the pattern
	var bestLetter = rune(0)
	var bestLower = rune(0)
	var bestLetterIdx = 0
	var bestLetterScore = 0
	var formattedStr = ""

	var matchedIndices []int

	// Loop over strings
	for strIdx != strLength {
		var patternChar = rune(0)
		if patternIdx != patternLength {
			patternChar = []rune(pattern)[patternIdx]
		}
		var strChar = []rune(str)[strIdx]

		var patternLower = rune(0)
		if patternChar != 0 {
			patternLower = unicode.ToLower(patternChar)
		}
		var strLower = unicode.ToLower(strChar)
		var strUpper = unicode.ToUpper(strChar)

		var nextMatch = patternChar != 0 && patternLower == strLower
		var rematch = bestLetter != 0 && bestLower == strLower

		var advanced = nextMatch && bestLetter != 0
		var patternRepeat = bestLetter != 0 && patternChar != 0 && bestLower == patternLower
		if advanced || patternRepeat {
			score += bestLetterScore
			matchedIndices = append(matchedIndices, bestLetterIdx)
			bestLetter = 0
			bestLower = 0
			bestLetterIdx = 0
			bestLetterScore = 0
		}

		if nextMatch || rematch {
			var newScore = 0

			// Apply penalty for each letter before the first pattern match
			// Note: std::max because penalties are negative values. So max is smallest penalty.
			if patternIdx == 0 {
				var penalty = max(strIdx*leading_letter_penalty, max_leading_letter_penalty)
				score += penalty
			}

			// Apply bonus for consecutive bonuses
			if prevMatched {
				newScore += adjacency_bonus
			}

			// Apply bonus for matches after a separator
			if prevSeparator {
				newScore += separator_bonus
			}

			// Apply bonus across camel case boundaries. Includes "clever" isLetter check.
			if prevLower && strChar == strUpper && strLower != strUpper {
				newScore += camel_bonus
			}

			// Update patter index IFF the next pattern letter was matched
			if nextMatch {
				patternIdx += 1
			}

			// Update best letter in str which may be for a "next" letter or a "rematch"
			if newScore >= bestLetterScore {

				// Apply penalty for now skipped letter
				if bestLetter != 0 {
					score += unmatched_letter_penalty
				}

				bestLetter = strChar
				bestLower = unicode.ToLower(bestLetter)
				bestLetterIdx = strIdx
				bestLetterScore = newScore
			}

			prevMatched = true
		} else {
			score += unmatched_letter_penalty
			prevMatched = false
		}

		// Includes "clever" isLetter check.
		prevLower = strChar == strLower && strLower != strUpper
		prevSeparator = strChar == '_' || strChar == ' '

		strIdx += 1
	}

	// Apply score for last match
	if bestLetter != 0 {
		score += bestLetterScore
		matchedIndices = append(matchedIndices, bestLetterIdx)
	}

	// Finish out formatted string after last pattern matched
	// Build formated string based on matched letters
	var lastIdx = 0
	rs := []rune(str)
	for _, idx := range matchedIndices {
		formattedStr += string(rs[lastIdx:idx-lastIdx]) + "<b>" + string(rs[idx]) + "</b>"
		lastIdx = idx + 1
	}
	formattedStr += string(rs[lastIdx:])

	var matched = patternIdx == patternLength
	return matched, score, formattedStr
}
