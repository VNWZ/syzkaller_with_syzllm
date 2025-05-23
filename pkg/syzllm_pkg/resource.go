package syzllm_pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// todo: impl ParseResource
func ParseResource(syzllmCall string, calls []string, insertPosition int) []string {
	calls[insertPosition] = syzllmCall

	// no res tag to parse but has res provider
	if !hasResTag(syzllmCall) {
		if countResProvider(syzllmCall) > 0 {
			insertCallNoResTag(calls, insertPosition, syzllmCall)
		}
		return calls
	}

	// todo: update nested res
	return calls // temp
}

func hasResTag(call string) bool {
	return hasNormalResTag(call) || hasPipeTag(call)
}

const (
	RPrefix    = "@RSTART@"
	RSuffix    = "@REND@"
	PIPEPrefix = "@PIPESTART@"
	PIPESuffix = "@PIPEEND@"
)

func hasNormalResTag(call string) bool {
	if strings.Contains(call, RPrefix) && strings.Contains(call, RSuffix) {
		return true
	}
	return false
}

func hasPipeTag(call string) bool {
	if strings.Contains(call, PIPEPrefix) && strings.Contains(call, PIPESuffix) {
		return true
	}
	return false
}

func assertNoResToParse(call string) {
	if hasResTag(call) {
		panic("Count Res Provider Must be in a cal without any res need to parse\nCall: " + call)
	}
}

func countResProvider(call string) int {
	assertNoResToParse(call)
	return countNormalResProvider(call) + countMultiResProvider(call)
}

var normalResProviderPattern = regexp.MustCompile(`^r(\d+) = `)

func countNormalResProvider(call string) int {
	assertNoResToParse(call)
	// return 0 or 1 since each call can provide up to 1 normal res
	if normalResProviderPattern.MatchString(call) {
		return 1
	}
	return 0
}

// todo: use combined regex
var multiResProviderPattern = regexp.MustCompile(`<r(\d+)=>`)

func countMultiResProvider(call string) int {
	assertNoResToParse(call)
	return len(multiResProviderPattern.FindAllString(call, -1))
}

// insert 1 new call at a time
func insertCallNoResTag(calls []string, insertPosition int, call string) {
	resProviderCountBeforeInsertPosition := 0
	for i, c := range calls {
		if i >= insertPosition {
			break
		}
		resProviderCountBeforeInsertPosition += countResProvider(c)
	}

	resProviderCountInNewCall := countResProvider(call)

	// update newcall res num
	calls[insertPosition] = updateResNum(call, resProviderCountBeforeInsertPosition)
	// update remain calls res num
	updatedRemainCalls := updateRemainCallsResNum(calls[insertPosition+1:], resProviderCountInNewCall, resProviderCountBeforeInsertPosition)
	for i, c := range updatedRemainCalls {
		calls[insertPosition+1+i] = c
	}
}

func updateRemainCallsResNum(remainCalls []string, offset int, resProviderCountBeforeInsertPosition int) []string {
	// Map to store original tag number to new number (e.g., "6" -> "6+offset")
	tagMap := make(map[string]string)

	result := make([]string, len(remainCalls))
	copy(result, remainCalls)

	// Process each string
	for i, s := range result {
		// Find all res prov in the current string
		matches := resProviderPattern.FindAllStringSubmatchIndex(s, -1)
		if len(matches) == 0 {
			continue
		}

		// Build the updated string
		var builder strings.Builder
		lastIndex := 0
		for _, match := range matches {
			// Write the part before the match
			builder.WriteString(s[lastIndex:match[0]])

			// Get the full matched res prov
			tag := s[match[0]:match[1]]
			var numStr string
			if match[2] != -1 {
				// rX = tag (group 1)
				numStr = s[match[2]:match[3]]
			} else {
				// <rX=> tag (group 2)
				numStr = s[match[4]:match[5]]
			}

			// Check if tag number is greater than previousCount
			num, _ := strconv.Atoi(numStr)
			if num >= resProviderCountBeforeInsertPosition {
				// Update tag with offset
				newNum := num + offset
				newNumStr := strconv.Itoa(newNum)
				tagMap[numStr] = newNumStr

				// Write the updated tag
				if strings.HasPrefix(tag, "<r") {
					builder.WriteString(fmt.Sprintf("<r%s=>", newNumStr))
				} else {
					builder.WriteString(fmt.Sprintf("r%s = ", newNumStr))
				}
			} else {
				// Keep the original tag
				builder.WriteString(tag)
			}

			lastIndex = match[1]
		}

		// Append the remaining part of the string
		builder.WriteString(s[lastIndex:])
		result[i] = builder.String()

		resReferencesFormats := []string{"r%s,", "r%s)", "r%s}"}
		// Update the same res prov numbers in subsequent strings
		for j := i + 1; j < len(result); j++ {
			for oldNum, newNum := range tagMap {
				for _, resRef := range resReferencesFormats {
					// Replace rX = (at start)
					oldTagStart := fmt.Sprintf(resRef, regexp.QuoteMeta(oldNum))
					newTagStart := fmt.Sprintf(resRef, newNum)
					result[j] = regexp.MustCompile(oldTagStart).ReplaceAllString(result[j], newTagStart)

					// Replace <rX=> (anywhere)
					oldTag := fmt.Sprintf(resRef, regexp.QuoteMeta(oldNum))
					newTag := fmt.Sprintf(resRef, newNum)
					result[j] = strings.ReplaceAll(result[j], oldTag, newTag)
				}
			}
		}
	}

	return result
}

var resProviderPattern = regexp.MustCompile(`^r(\d+) = |<r(\d+)=>`)

func updateResNum(call string, previousCount int) string {
	matches := resProviderPattern.FindAllStringSubmatchIndex(call, -1)
	if len(matches) == 0 {
		return call
	}

	var result strings.Builder
	lastIndex := 0
	currentCount := previousCount

	for _, match := range matches {
		result.WriteString(call[lastIndex:match[0]])

		tag := call[match[0]:match[1]]

		if strings.HasPrefix(tag, "<r") {
			// Handle <rZ=> tag
			result.WriteString(fmt.Sprintf("<r%d=>", currentCount))
		} else if strings.HasSuffix(tag, " = ") {
			// Handle rX = tag
			result.WriteString(fmt.Sprintf("r%d = ", currentCount))
		}

		currentCount++
		lastIndex = match[1]
	}

	result.WriteString(call[lastIndex:])

	return result.String()
}
