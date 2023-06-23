package documentize

import "strings"

// stringToWords takes an input string s and converts it to words.
func stringToWords(s string) string {
	stripModule := strings.Replace(strings.Replace(s, "module.", "", -1), ".", " ", -1)

	// handling the case where there are '/' in the string
	if strings.Contains(stripModule, "/") {
		stripModule = stripModule[strings.LastIndex(stripModule, "/")+1:]
	}

	edgeCaseUnderScoreAndDash := strings.ReplaceAll(strings.ReplaceAll(stripModule, "_-", " "), "-_", " ")

	splitUnderScore := strings.ReplaceAll(edgeCaseUnderScoreAndDash, "_", " ")

	splitDash := strings.ReplaceAll(splitUnderScore, "-", " ")

	return splitDash
}
