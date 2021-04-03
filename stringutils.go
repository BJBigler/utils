package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//Tokenize returns a slice of unique "words" of *minLength*
//from the given string *args". The tokenization proceeds
//in two ways: 1) words into "phrases" and 2) words into
//"letters". The phrases are full words recombined into
//searchable phrases; e.g., "Our Lady of Los Angeles" becomes
//our, our lady, ..., Los, Los Angeles.
//The letters are strings starting from the beginning of
//each term. For example, if the terms
//supplied are "the thames", and the minimum length was 2,
//the result would be {th, the, tha, tham, thame, thames}
func Tokenize(minLength int, terms ...string) []string {

	//Use a map to enforce term uniqueness
	uniques := make(map[string]string)

	for _, t := range terms {

		if len(t) < minLength {
			continue
		}

		//remove double spaces
		t = strings.Replace(t, "  ", " ", -1)

		//1) Add words as phrases
		phrases := TokenizeWords(t)
		for _, w := range phrases {
			if w == "" {
				continue
			}

			uniques[w] = w
		}

		//2) Add words as fragments
		runeTerms := []rune(t)

		//Split the words into letters
		for j := 0; j < len(runeTerms)+1; j++ {

			fragment := runeTerms[0:j]

			if len(fragment) < minLength {
				continue
			}

			uniques[string(fragment)] = string(fragment)
		}
	}

	result := []string{}

	for v := range uniques {
		result = append(result, v)
	}

	return result
}

//TokenizeWords splits a phrase into words and returns
//phrases 1 .. n, 2 .. n. For example, "General Smith of Los Angeles"
//would return:
//General, General Smith, ..., Smith, Smith of, ..., of, of Los, ...,
//Los, Los Angeles
func TokenizeWords(term string) (result []string) {
	term = strings.ToLower(term)

	terms := strings.Fields(term)

	termLen := len(terms)

	//This func does all the work,
	//starting at 'start' -- i.e., word 0, word 1, etc. --
	//and grabbing that word and those beyond in turn to
	//get phrases
	addterms := func(start int) {
		for i := start; i <= termLen; i++ {

			phrase := strings.Join(terms[start:i], " ")
			if phrase == "" {
				continue
			}

			result = append(result, phrase)
		}
	}

	for start := range terms {
		addterms(start)
	}

	return result
}

const letterBytes = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//GenerateRandomAlphaNumeric ...
func GenerateRandomAlphaNumeric(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

//ReplaceAccents replaces accented characters with close eqivalents
func ReplaceAccents(input string) string {

	input = strings.Replace(input, "æ", "ae", -1)

	// isMn := func(r rune) bool {
	// 	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	// }

	b := make([]byte, len(input))

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	//t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, err := t.Transform(b, []byte(input), true)

	if err != nil {
		return input
	}

	return string(b)

	// t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	// s, _, err := transform.String(t, input)

	// if err != nil {

	// }

}

//Conjoin joins a string with commas and the conjunction,
//e.g., ["bill", "sue", "jill"] => "bill, sue, and jill"
func Conjoin(val []string, conjunction string) string {

	lenVal := len(val)

	switch lenVal {
	case 0:
		return ""
	case 1:
		return val[0]
	case 2:
		return strings.Join(val, fmt.Sprintf(" %s ", conjunction))
	default:
		return fmt.Sprintf("%s, %s %s", strings.Join(val[0:lenVal-1], ", "), conjunction, val[lenVal-1])
	}
}

//PrepCookieDomain removes http(s) from the institution's root URL
//to function as a cookie domain value
//e.g., https://princeton.seminars.app -> princeton.seminars.app
func PrepCookieDomain(val string) string {
	val = strings.Replace(val, "https://", "", -1)
	val = strings.Replace(val, "http://", "", -1)
	// if strings.HasPrefix(val, ".") {
	// 	val = val[1:]
	// }
	return val
}

//MakeFullname concats first and last name
func MakeFullname(firstname, lastname, alternate string) string {
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	alternate = strings.TrimSpace(alternate)

	if firstname == "" && lastname == "" {
		return ""
	}

	if alternate != "" {
		alternate = fmt.Sprintf(" (%s)", alternate)
	}

	return fmt.Sprintf("%s %s%s", firstname, lastname, alternate)
}
