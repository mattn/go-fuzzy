package fuzzy

import (
	"unicode"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Match(pattern, str string, matchedIndices *[]int) (bool, int) {
	// Score consts
	const (
		adjacencyBonus          = 5  // bonus for adjacent matches
		separatorBonus          = 10 // bonus if match occurs after a separator
		camelBonus              = 10 // bonus if match is uppercase and prev is lower
		leadingLetterPenalty    = -3 // penalty applied for every letter in str before the first match
		maxLeadingLetterPenalty = -9 // maximum penalty for leading letters
		unmatchedLetterPenalty  = -1 // penalty for every letter that doesn't matter
	)

	// Loop variables
	var score = 0
	var rp = []rune(pattern)
	var rs = []rune(str)
	var ip = 0
	var lp = len(rp)
	var is = 0
	var ls = len(rs)
	var prevMatched = false
	var prevLower = false
	var prevSeparator = true // true so if first letter match gets separator bonus

	// Use "best" matched letter if multiple string letters match the pattern
	var bestLetter = rune(0)
	var bestLower = rune(0)
	var bestLetterIdx = 0
	var bestLetterScore = 0

	// Loop over strings
	for is < ls {
		var cp = rune(0)
		if ip != lp {
			cp = rp[ip]
		}
		var cs = rs[is]

		var patternLower = rune(0)
		if cp != 0 {
			patternLower = unicode.ToLower(cp)
		}
		var strLower = unicode.ToLower(cs)
		var strUpper = unicode.ToUpper(cs)

		var nextMatch = cp != 0 && patternLower == strLower
		var rematch = bestLetter != 0 && bestLower == strLower

		var advanced = nextMatch && bestLetter != 0
		var patternRepeat = bestLetter != 0 && cp != 0 && bestLower == patternLower
		if advanced || patternRepeat {
			score += bestLetterScore
			if matchedIndices != nil {
				*matchedIndices = append(*matchedIndices, bestLetterIdx)
			}
			bestLetter = 0
			bestLower = 0
			bestLetterIdx = 0
			bestLetterScore = 0
		}

		if nextMatch || rematch {
			var newScore = 0

			// Apply penalty for each letter before the first pattern match
			// Note: std::max because penalties are negative values. So max is smallest penalty.
			if ip == 0 {
				score += max(is*leadingLetterPenalty, maxLeadingLetterPenalty)
			}

			// Apply bonus for consecutive bonuses
			if prevMatched {
				newScore += adjacencyBonus
			}

			// Apply bonus for matches after a separator
			if prevSeparator {
				newScore += separatorBonus
			}

			// Apply bonus across camel case boundaries. Includes "clever" isLetter check.
			if prevLower && cs == strUpper && strLower != strUpper {
				newScore += camelBonus
			}

			// Update patter index IFF the next pattern letter was matched
			if nextMatch {
				ip += 1
			}

			// Update best letter in str which may be for a "next" letter or a "rematch"
			if newScore >= bestLetterScore {

				// Apply penalty for now skipped letter
				if bestLetter != 0 {
					score += unmatchedLetterPenalty
				}

				bestLetter = cs
				bestLower = unicode.ToLower(bestLetter)
				bestLetterIdx = is
				bestLetterScore = newScore
			}

			prevMatched = true
		} else {
			score += unmatchedLetterPenalty
			prevMatched = false
		}

		// Includes "clever" isLetter check.
		prevLower = cs == strLower && strLower != strUpper
		prevSeparator = cs == '_' || cs == ' '

		is += 1
	}

	// Apply score for last match
	if bestLetter != 0 {
		score += bestLetterScore
		if matchedIndices != nil {
			*matchedIndices = append(*matchedIndices, bestLetterIdx)
		}
	}

	return ip == lp, score
}
