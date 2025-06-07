package syzllm_pkg

import (
	"fmt"
	"github.com/google/syzkaller/pkg/log"
	"regexp"
	"strconv"
	"strings"
)

// calls[insertPosition] should be "[MASK]" at the beginning
func ParseResource(calls []string, insertPosition int, syzllmCall string) []string {
	// new call has no res tag to parse but has res provider
	if !hasResTag(syzllmCall) {
		if countResProvider(syzllmCall) > 0 {
			insertCallNoResTag(calls, insertPosition, syzllmCall)
		} else {
			calls[insertPosition] = syzllmCall
		}
		return calls
	}

	// 1. parsedSyzllmCalls
	parsedSyzllmCalls := parseResTag(syzllmCall)
	// 2. check and insert one by one
	result := insertSyzllmCalls(calls, insertPosition, parsedSyzllmCalls)

	return result
}

func parseResTag(syzllmCall string) []string {
	result := []string{}
	N := 0
	modified := strings.Builder{}
	pos := 0

	for pos < len(syzllmCall) {
		rStart := strings.Index(syzllmCall[pos:], RPrefix)
		if rStart != -1 {
			rStart += pos
		}
		pipeStart := strings.Index(syzllmCall[pos:], PIPEPrefix)
		if pipeStart != -1 {
			pipeStart += pos
		}

		if rStart == -1 && pipeStart == -1 {
			break
		}

		// determine which type is first
		var startPos int
		var tagType string
		if rStart != -1 && (pipeStart == -1 || rStart < pipeStart) {
			startPos = rStart
			tagType = "res"
		} else {
			startPos = pipeStart
			tagType = "pipe"
		}

		var startTag, endTag string
		if tagType == "res" {
			startTag = RPrefix
			endTag = RSuffix
		} else {
			startTag = PIPEPrefix
			endTag = PIPESuffix
		}

		endPos := strings.Index(syzllmCall[startPos+len(startTag):], endTag)
		if endPos == -1 {
			panic("No matching end tag")
		}
		endPos += startPos + len(startTag)

		substring := syzllmCall[startPos+len(startTag) : endPos]

		// start process the first res tag
		if tagType == "res" {
			result = append(result, fmt.Sprintf("r%d = %s", N, substring))
			modified.WriteString(syzllmCall[pos:startPos])
			modified.WriteString(fmt.Sprintf("r%d", N))
			N++
		} else { // pipe tag
			re := regexp.MustCompile(`<r\d+=>`)
			matches := re.FindAllStringIndex(substring, -1)
			k := len(matches)
			if k == 0 {
				panic("No placeholders in pipe tag")
			}
			firstN := N
			updatedSubstring := substring
			for i := 0; i < len(matches); i++ {
				matchStart := matches[i][0]
				matchEnd := matches[i][1]
				updatedSubstring = updatedSubstring[:matchStart] + fmt.Sprintf("<r%d=>", N) + updatedSubstring[matchEnd:]
				N++
			}
			result = append(result, updatedSubstring)
			modified.WriteString(syzllmCall[pos:startPos])
			modified.WriteString(fmt.Sprintf("r%d", firstN))
		}

		pos = endPos + len(endTag)
	}

	modified.WriteString(syzllmCall[pos:])
	result = append(result, modified.String())
	return result
}

func insertSyzllmCalls(calls []string, insertPos int, syzllmCalls []string) []string {
	result := make([]string, len(calls))
	copy(result, calls)

	cursor := insertPos

	tagMap := make(map[string]string)

	syzllmCallsSize := len(syzllmCalls)
	// if produced resource
	if syzllmCallsSize > 1 {
		for _, call := range syzllmCalls[0 : syzllmCallsSize-1] {
			resCount := hasCallName(result, 0, cursor, call)
			if resCount >= 0 {
				syzllmCallResNum := extractFirstResProvider(call)
				tagMap[syzllmCallResNum[1:]] = strconv.Itoa(resCount)
			} else {
				result, _ = enlargeSlice(result, cursor)
				insertCallNoResTag(result, cursor, call)
				oldResNum := extractFirstResProvider(call)
				updatedResNum := extractFirstResProvider(result[cursor])
				tagMap[oldResNum[1:]] = updatedResNum[1:]

				cursor += 1
			}
		}
	}

	// update res by tagMap
	sortedKeys := SortedMapKeysDesc(tagMap)
	// Update the same res prov numbers in subsequent strings
	for _, k := range sortedKeys {
		oldResNum := "r" + k
		updatedResNum := "r" + tagMap[k]
		syzllmCalls[syzllmCallsSize-1] = strings.Replace(syzllmCalls[syzllmCallsSize-1], oldResNum+",", updatedResNum+",", -1)
		syzllmCalls[syzllmCallsSize-1] = strings.Replace(syzllmCalls[syzllmCallsSize-1], oldResNum+")", updatedResNum+")", -1)
		syzllmCalls[syzllmCallsSize-1] = strings.Replace(syzllmCalls[syzllmCallsSize-1], oldResNum+"}", updatedResNum+"}", -1)
	}

	insertCallNoResTag(result, cursor, syzllmCalls[syzllmCallsSize-1])
	return result
}

func hasCallName(calls []string, start int, end int, syzllmCall string) int {
	for i := start; i < end; i++ {
		if extractCallName(calls[i]) == extractCallName(syzllmCall) {
			resProvider := extractFirstResProvider(calls[i])
			if len(resProvider) <= 1 {
				log.Errorf("Failed to extract res prov from %s, syzllm call %s", calls[i], syzllmCall)
			}
			ret, _ := strconv.Atoi(resProvider[1:])
			return ret
		}
	}
	return -1
}

func extractCallName(call string) string {
	// Find the position of the first '$' or '('
	end := len(call)
	for i, c := range call {
		if c == '$' || c == '(' {
			end = i
			break
		}
	}

	// Extract the substring up to the found position
	result := call[:end]

	// If the string contains '=', trim up to and including '='
	if idx := strings.Index(result, "="); idx != -1 {
		result = strings.TrimSpace(result[idx+1:])
	}

	return result
}

func extractFirstResProvider(call string) string {
	// Check if the string starts with a resource tag (e.g., "rX =")
	if idx := strings.Index(call, "="); idx != -1 && idx > 1 && call[0] == 'r' {
		return strings.TrimSpace(call[:idx])
	}

	// Look for resource tags in the format <rX=...> within the string
	start := strings.Index(call, "<r")
	if start == -1 {
		return ""
	}
	end := strings.Index(call[start:], ">")
	if end == -1 {
		return ""
	}
	end += start

	// Extract the first resource tag (e.g., "r1" from "<r1=...")
	tag := call[start+1 : end]
	if idx := strings.Index(tag, "="); idx != -1 {
		return tag[:idx]
	}
	return tag
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
		panic("Count Res Provider Must be in a call without any res need to parse\nCall: " + call)
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

var multiResProviderPattern = regexp.MustCompile(`<r(\d+)=>`)

func countMultiResProvider(call string) int {
	assertNoResToParse(call)
	return len(multiResProviderPattern.FindAllString(call, -1))
}

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
	calls[insertPosition] = updateResProvider(call, resProviderCountBeforeInsertPosition)
	// update remain calls res num
	updatedRemainCalls := updateRemainCallsResNum(calls[insertPosition+1:], resProviderCountInNewCall, resProviderCountBeforeInsertPosition)
	for i, c := range updatedRemainCalls {
		calls[insertPosition+1+i] = c
	}
}

// 1. assign new res num for Prov.
// 2. update reference num by map
func updateRemainCallsResNum(remainCalls []string, resProviderCountAtInsertPosition int, resProviderCountBeforeInsertPosition int) []string {
	count := resProviderCountBeforeInsertPosition + resProviderCountAtInsertPosition

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
				newNum := count
				count += 1
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
	}

	// keys must be sorted avoiding overwriting
	// e.g., r2->r3, r3->r4; replacing call(r2, r3), 1st loop call(r3, r4) -> 2nd loop call(r4, r4)
	sortedKeys := SortedMapKeysDesc(tagMap)

	resReferencesFormats := []string{"r%s,", "r%s)", "r%s}"}
	// Update the same res prov numbers in subsequent strings
	for i := 0; i < len(result); i++ {
		for _, k := range sortedKeys {
			oldNum := k
			newNum := tagMap[k]
			for _, resRef := range resReferencesFormats {
				// Replace rX = (at start)
				oldTagStart := fmt.Sprintf(resRef, regexp.QuoteMeta(oldNum))
				newTagStart := fmt.Sprintf(resRef, newNum)
				result[i] = strings.ReplaceAll(result[i], oldTagStart, newTagStart)

				// Replace <rX=> (anywhere)
				oldTag := fmt.Sprintf(resRef, regexp.QuoteMeta(oldNum))
				newTag := fmt.Sprintf(resRef, newNum)
				result[i] = strings.ReplaceAll(result[i], oldTag, newTag)
			}
		}
	}

	return result
}

var resProviderPattern = regexp.MustCompile(`^r(\d+) = |<r(\d+)=>`)

func updateResProvider(call string, previousCount int) string {
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
