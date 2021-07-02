package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/mail"
	"net/url"
	"unicode"

	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

//Alpha returns A-Z, capitalized
func Alpha() []string {
	return []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
}

//ParseInt ...
func ParseInt(candidate string, defaultResult int) int {

	result, err := strconv.Atoi(candidate)

	if err != nil {
		return defaultResult
	}

	return result
}

//ParseInt32 ...
func ParseInt32(candidate string, defaultResult int32) int32 {

	result, err := strconv.ParseInt(candidate, 10, 32)

	if err != nil {
		return defaultResult
	}

	return int32(result)
}

//ParseInt64 ...
func ParseInt64(candidate string, defaultResult int64) int64 {

	candidate = strings.Replace(candidate, "–", "-", -1)
	candidate = strings.Replace(candidate, "—", "-", -1)

	result, err := strconv.ParseInt(candidate, 10, 64)

	if err != nil {
		return defaultResult
	}

	return result
}

//ParseInt64Err returns an error if parsing fails
func ParseInt64Err(candidate string) (int64, error) {

	candidate = strings.Replace(candidate, "–", "-", -1)
	candidate = strings.Replace(candidate, "—", "-", -1)

	result, err := strconv.ParseInt(candidate, 10, 64)

	if err != nil {
		return 0, err
	}

	return result, nil
}

//ParseFloat32 ...
func ParseFloat32(candidate string) float32 {

	candidate = strings.TrimSpace(candidate)

	if candidate == "-" || candidate == "" {
		return 0
	}

	result, err := strconv.ParseFloat(candidate, 32)

	if err != nil {
		return 0
	}

	return float32(result)
}

//ParseFloat64 ...
func ParseFloat64(candidate string) float64 {

	candidate = strings.TrimSpace(candidate)

	if candidate == "-" || candidate == "" {
		return 0
	}

	result, err := strconv.ParseFloat(candidate, 64)

	if err != nil {
		return 0
	}

	return result
}

//ParseDateMulti tries to parse the provided value using
//three formats:
//1) yyyy-mm-dd
//2) yyyy-mm-dd hh:mm:ss
//3)  yyyy-mm-dd hh:mm
//If all these produce errors, time.Time{} is returned
func ParseDateMulti(candidate string, location *time.Location) (result time.Time) {

	result, err := ParseDate([]byte(candidate), location)
	if err != nil {
		//yyyy-mm-dd hh:mm:ss
		result, err = ParseDateTime4(candidate, location)
		if err != nil {
			//yyyy-mm-dd hh:mm
			result, err = ParseDateTime5(candidate, location)
			if err != nil {
				result = time.Time{}
			}
		}
	}

	return
}

//IntToString ...
func IntToString(val int) string {

	return strconv.Itoa(val)
}

//ParseBool ...
func ParseBool(val string) bool {

	val = strings.ToLower(val)
	val = strings.TrimSpace(val)

	if val == "true" || val == "t" || val == "1" || val == "yes" {
		return true
	}

	return false

}

//Int64ToString ...
func Int64ToString(val int64) string {

	return strconv.FormatInt(val, 10)
}

//OnToBool ...
func OnToBool(val string) bool {

	val = strings.ToLower(val)

	if val == "on" || val == "true" || val == "1" || val == "yes" {
		return true
	}

	return false
}

//OnToInt returns 1 if on, 0 otherwise
func OnToInt(val string) int64 {
	if strings.ToLower(val) == "on" {
		return 1
	}

	return 0
}

//StringToBool returns true if "true", "1", "on", or "yes"
func StringToBool(val string) bool {
	val = strings.ToLower(val)

	val = strings.TrimSpace(val)

	if val == "true" {
		return true
	}

	if val == "1" {
		return true
	}

	if val == "on" {
		return true
	}

	if val == "yes" {
		return true
	}

	return false
}

//ValToBool returns true if val is "on",
//"true", or "1"; false for all others.
//val is converted to lower before checking.
func ValToBool(val string) bool {
	val = strings.ToLower(val)

	if val == "true" || val == "on" || val == "1" {
		return true
	}

	return false
}

//RemovePunctuation ...
func RemovePunctuation(input string) string {

	var out strings.Builder

	for _, l := range input {
		//Not a letter and not a space
		if !unicode.IsLetter(l) && !unicode.IsSpace(l) {
			out.WriteString(" ")
			continue
		}

		out.WriteString(string(l))
	}

	return out.String()

	// punc := "\"!`#$%&'()*+,-./:;?@[\\]^_{|}~"

	// out := input

	// if strings.ContainsAny(input, punc) {
	// 	out = ""

	// 	for _, c := range input {
	// 		for _, p := range punc {
	// 			if c == p {
	// 				out = " "
	// 				continue
	// 			}
	// 			out += string(c)
	// 		}
	// 	}
	// }

	// return out

}

//RemoveSpaces trims spaces around a string, and removes two+ spaces inside string
func RemoveSpaces(input string) string {

	reLeadcloseWhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	input = reLeadcloseWhtsp.ReplaceAllString(input, "")
	input = reInsideWhtsp.ReplaceAllString(input, " ")

	return input

}

//IsUpperCaseWord checks if a word has all uppercase characters
func IsUpperCaseWord(word string) bool {

	for _, char := range word {
		if unicode.IsLower(char) {
			return false
		}
	}

	return true
}

//ToTitleCase returns the input string in Title Case
func ToTitleCase(ctx context.Context, phrase string) string {
	if IsUpperCaseWord(phrase) {
		//The entire name is in uppercase.
		//Convert it to title case. This requires first setting all characters
		//to lower case, and then changing it to Title
		phrase = strings.ToLower(phrase)
		phrase = strings.Title(phrase)
		return phrase
	}
	return phrase
}

//ToInt64ForStorage multiplies the input number by 10^precision
//and returns an int64 for Datastore persist
func ToInt64ForStorage(number decimal.Decimal, precision int32) int64 {

	decimalMultiplier := decimal.New(1, precision)
	number = number.Mul(decimalMultiplier)

	number = number.Round(0)
	return number.IntPart()
}

//ToInt64Precision10ForStorage multiplies the input number by 10^10
//and returns an int64 for Datastore persist
func ToInt64Precision10ForStorage(number decimal.Decimal) int64 {

	decimalMultiplier := decimal.New(1, 10)
	number = number.Mul(decimalMultiplier)

	number = number.Round(0)
	return number.IntPart()
}

//FormatCommas adds commas to a number
func FormatCommas(num string) string {

	re := regexp.MustCompile(`(\d+)(\d{3})`)
	for {
		formatted := re.ReplaceAllString(num, "$1,$2")
		if formatted == num {
			return formatted
		}
		num = formatted
	}
}

//PrepForDecimalization preps string for decimalization.
func PrepForDecimalization(num string) string {
	num = strings.Replace(num, ",", "", -1)
	num = strings.TrimSpace(num)
	return num
}

//GenerateRandom returns a 7-digit pseudorandom number
func GenerateRandom() int64 {
	rand.Seed(time.Now().Unix())
	min := 1000000
	max := 9999999
	return int64(rand.Intn(max-min) + min)
}

//GenerateRandom returns a pseudorandom whole number in the provided range
func GenerateRandomFromRange(min, max int) int64 {
	rand.Seed(time.Now().Unix())
	return int64(rand.Intn(max-min) + min)
}

//GetURL goes to the given URL and returns whatever html/string at the address.
// func GetURL(ctx context.Context, url string) string {

// 	http.Get()
// 	client := urlfetch.Client(ctx)
// 	rsp, err := client.Get(url)

// 	if err != nil {
// 		//Logger(ctx, err)
// 		return ""
// 	}

// 	defer rsp.Body.Close()
// 	body, err := ioutil.ReadAll(rsp.Body)

// 	if err != nil {
// 		Log(err)
// 	}

// 	return string(body)
// }

//GetIDFromURL returns ID from a URL fragment, e.g.
// the URL http://abc.com/seminar/1234 would return 1234.
//Note that the ID has to be numerical
func GetIDFromURL(url *url.URL) int64 {
	lastPart := GetStringIDFromURL(url)
	return ParseInt64(lastPart, 0)
}

//GetStringIDFromURL gets the last text past the final
//slash in a URL
func GetStringIDFromURL(url *url.URL) string {
	parts := strings.Split(url.Path, "/")
	lastPart := len(parts) - 1
	return parts[lastPart]
}

//GetLastPart gets the last text past the final
//delimiter
func GetLastPart(path, delimiter string) string {
	parts := strings.Split(path, delimiter)
	lastPart := len(parts) - 1
	return parts[lastPart]
}

//Int64ToCSV converts a list of int64s
//to CSV
func Int64ToCSV(vals []int64) string {
	var IDs []string
	for _, i := range vals {
		IDs = append(IDs, fmt.Sprintf("%v", i))
	}

	return strings.Join(IDs, ", ")
}

//SplitSubN splits a string by length
func SplitSubN(s string, numberOfCharactersPerLine int) []string {
	n := numberOfCharactersPerLine

	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

//GetConjunction is used to make a list human-readable with a comma and/or the
//supplied conjunction. For example, if the slice ["jack","jane","mary"]
//should be formatted "jack, jane, and mary" or jack, jane, or mary" --
//GetConjuction will supply the two commas and an 'and' or 'or' as fed to it.
func GetConjunction(cntTotal int, current int, conjunction string) string {

	if cntTotal == 2 && current == 0 {
		return fmt.Sprintf(" %v ", conjunction) //Judy or Joey
	}
	if cntTotal > 2 {
		if current < cntTotal-2 { //items 0 to n, but not penultimate
			return ", "
		}
		if current == cntTotal-2 {
			return fmt.Sprintf(", %v ", conjunction) //penultimate item
		}
	}
	return ""
}

//IsNumeric removes commas, trims spaces, and then attempts to parse as float.
func IsNumeric(s string) bool {

	s = strings.Replace(s, ",", "", -1)
	s = strings.TrimSpace(s)

	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

//ValidateEmail checks email for ...
//1) More than 0 length;
//2) Contains the provided suffix (e.g., princeton.edu); and
//3) Properly parsed by Go's built-in mail.ParseAddress;
func ValidateEmail(s, domain string) (bool, error) {
	if len(s) == 0 {
		return false, fmt.Errorf("email is blank")
	}

	_, err := mail.ParseAddress(s)
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(domain, "@") {
		domain = "@" + domain
	}
	if !strings.HasSuffix(s, domain) {
		return false, fmt.Errorf("email does not end with %s", domain)
	}

	return true, nil
}

//MakeSliceFromRange returns a slice of integers in the range 'from'-'to' inclusive
func MakeSliceFromRange(fromVal, toVal int) (result []int) {
	for i := fromVal; i < toVal; i++ {
		result = append(result, i)
	}
	result = append(result, toVal)
	return
}

//SumFloat64 ...
func SumFloat64(vals []float64) (result float64) {
	for _, v := range vals {
		result += v
	}
	return result
}

//SumInt64 ...
func SumInt64(vals []int64) (result int64) {
	for _, v := range vals {
		result += v
	}
	return result
}

//AlphaIndex returns the 0-based index
//where a letter appears in the English alphabet.
//For example, "E" would return 5.
func AlphaIndex(ltr string) int {
	for i, l := range Alpha() {
		if l != strings.ToLower(ltr) {
			continue
		}
		return i
	}
	return -1
}

//EmojiClockFace returns the clock face for
//times from 8:00 through 21:30 in half-hour
//increments. The default is 2:00
func EmojiClockFace(t time.Time) string {

	clockEmoji := "🕑"

	switch t.Format("15:04") {
	case "08:00":
		clockEmoji = "🕗"
	case "08:30":
		clockEmoji = "🕣"
	case "09:00":
		clockEmoji = "🕘"
	case "09:30":
		clockEmoji = "🕤"
	case "10:00":
		clockEmoji = "🕙"
	case "10:30":
		clockEmoji = "🕥"
	case "11:00":
		clockEmoji = "🕚"
	case "11:30":
		clockEmoji = "🕦"
	case "12:00":
		clockEmoji = "🕛"
	case "12:30":
		clockEmoji = "🕧"
	case "13:00":
		clockEmoji = "🕐"
	case "13:30":
		clockEmoji = "🕜"
	case "14:00":
		clockEmoji = "🕑"
	case "14:30":
		clockEmoji = "🕝"
	case "15:00":
		clockEmoji = "🕒"
	case "15:30":
		clockEmoji = "🕞"
	case "16:00":
		clockEmoji = "🕓"
	case "16:30":
		clockEmoji = "🕟"
	case "17:00":
		clockEmoji = "🕔"
	case "17:30":
		clockEmoji = "🕠"
	case "18:00":
		clockEmoji = "🕕"
	case "18:30":
		clockEmoji = "🕡"
	case "19:00":
		clockEmoji = "🕖"
	case "19:30":
		clockEmoji = "🕢"
	case "20:00":
		clockEmoji = "🕗"
	case "20:30":
		clockEmoji = "🕣"
	case "21:00":
		clockEmoji = "🕘"
	case "21:30":
		clockEmoji = "🕤"
	}

	return clockEmoji
}

//ZeroPad64 adds a zero in front of *val*
func ZeroPad64(val int64) string {
	return fmt.Sprintf("%02d", val)
}
