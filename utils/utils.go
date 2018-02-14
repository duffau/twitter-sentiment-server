package utils

import (
	"fmt"
	"regexp"
	"strings"
) 

var (
	eyes = `[8:=;]`
    nose = `['` + "`" + `\-]?`
	flags = `(?ms)`

	url = regexp.MustCompile(flags + `https?:\/\/\S+\b|www\.(\w+\.)+\S*`)
	slash = regexp.MustCompile(flags + `/`)
	user = regexp.MustCompile(flags + `@\w+`)
	smile = regexp.MustCompile(flags + fmt.Sprintf(`%s%s[)dD]+|[)dD]+%s%s`, eyes, nose, nose, eyes))
    lolface = regexp.MustCompile(flags + fmt.Sprintf(`%s%sp+`, eyes, nose))
    sadface = regexp.MustCompile(flags + fmt.Sprintf(`%s%s\(+|\)+%s%s`, eyes, nose, nose, eyes))
    neutralface = regexp.MustCompile(flags + fmt.Sprintf(`%s%s[\/|l*]`, eyes, nose))
    heart = regexp.MustCompile(flags + `<3`)
    number = regexp.MustCompile(flags + `[-+]?[.\d]*[\d]+[:,.\d]*`)
    hastag = regexp.MustCompile(flags + `#\S+`)
    // repeat = regexp.MustCompile(flags + `([!?.]){2,}`)
    // elong = regexp.MustCompile(flags + `\b(\S*?)(.)$2{2,}\b`)
)


func hashtag(text string) string {
	var result string
	hashtagBody := text[1:]
	if strings.ToUpper(hashtagBody) == hashtagBody {
		result = fmt.Sprintf(`<hashtag> %s <allcaps>`, hashtagBody)
	} else {
		result = fmt.Sprintf(`<hashtag> %s`, hashtagBody)
	}
	return result
}

func allcaps(text string) string {
	return strings.ToLower(text) + ` <allcaps>`
}

func Tokenize(text string) string {
	text = url.ReplaceAllString(text, `<url>`)
    text = slash.ReplaceAllString(text,` / `)
    text = user.ReplaceAllString(text, `<user>`)
    text = smile.ReplaceAllString(text, `<smile>`)
    text = lolface.ReplaceAllString(text, `<lolface>`)
    text = sadface.ReplaceAllString(text, `<sadface>`)
    text = neutralface.ReplaceAllString(text, `<neutralface>`)
    text = heart.ReplaceAllString(text,`<heart>`)
    text = number.ReplaceAllString(text, `<number>`)
    text = hastag.ReplaceAllStringFunc(text, hashtag)
    // text = repeat.ReplaceAllString(text, `$1 <repeat>`)
    // text = elong.ReplaceAllString(text, `$1$2 <elong>`)

	return strings.ToLower(text)
}


