package util

import (
	"os"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewTextLogger(lvl log.Level, colors bool, out *os.File) *log.Logger {
	lg := log.New()

	formatter := &textFormatter{colors: colors}
	lg.SetFormatter(formatter)

	//lg.SetOutput(os.Stdout)
	lg.SetOutput(out)

	lg.SetLevel(lvl)

	return lg
}

type textFormatter struct {
	colors bool
}

func resolveTags(data map[string]interface{}, first bool) string {
	var line string
	var logFieldsType = reflect.TypeOf((*log.Fields)(nil)).Elem()
	for _, k := range sortedKeys(data) {
		d := data[k]
		if first {
			first = false
		} else {
			line += " "
		}
		if reflect.TypeOf(d) == logFieldsType {
			line += fmt.Sprintf("${lvlcolor}%+v={%+v${lvlcolor}}${RESET}", k, resolveTags(d.(log.Fields), true))
		} else {
			line += fmt.Sprintf("${lvlcolor}%+v=${BOLD}%+v${RESET}", k, d)
		}
	}
	return line
}

func (f *textFormatter) Format(e *log.Entry) ([]byte, error) {
	line := lvlTags[e.Level]
	line += timestampTag
	line += " " + msgTag
	line += "${justify}"
	line += resolveTags(e.Data, false)
	line = f.resolve(line, e)
	return []byte(line + "\n"), nil
}

var lvlColors = map[log.Level]string{
	log.PanicLevel: "${MAGENTA}${BOLD}",
	log.FatalLevel: "${MAGENTA}",
	log.ErrorLevel: "${RED}",
	log.WarnLevel:  "${YELLOW}",
	log.InfoLevel:  "${CYAN}",
	log.DebugLevel: "${GREEN}",
	log.TraceLevel: "${BLUE}",
}

var lvlTags = map[log.Level]string{
	log.PanicLevel: "${lvlcolor}[${lvl}]${RESET}",
	log.FatalLevel: "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
	log.ErrorLevel: "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
	log.WarnLevel:  "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
	log.InfoLevel:  "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
	log.DebugLevel: "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
	log.TraceLevel: "${lvlcolor}[${BOLD}${lvl}${RESET}${lvlcolor}]${RESET}",
}

const timestampTag = "${lvlcolor}[${timestamp}]${RESET}"

const msgTag = "${lvlcolor}${msg}${RESET}"

const justify = 75

const (
	blackC   = 30
	redC     = 31
	greenC   = 32
	yellowC  = 33
	blueC    = 34
	magentaC = 35
	cyanC    = 36
	whiteC   = 37
	resetC   = -1
)

const (
	colorSeq = "\033[%dm"
	boldSeq  = "\033[1m"
	resetSeq = "\033[0m"
)

func codeOf(c int) string {
	if c == resetC {
		return resetSeq
	}
	return fmt.Sprintf(colorSeq, c)
}

var colorCodes = map[string]string{
	"${BLACK}":   codeOf(blackC),
	"${RED}":     codeOf(redC),
	"${GREEN}":   codeOf(greenC),
	"${YELLOW}":  codeOf(yellowC),
	"${BLUE}":    codeOf(blueC),
	"${MAGENTA}": codeOf(magentaC),
	"${CYAN}":    codeOf(cyanC),
	"${WHITE}":   codeOf(whiteC),
	"${BOLD}":    boldSeq,
	"${RESET}":   resetSeq,
}

func lvlStr(l log.Level) string {
	return strings.ToUpper(l.String()[:1])
}

func timestampStr(t time.Time) string {
	s := t.Format("20060102-15:04:05")
	nanos := t.UnixNano() % 1000000000
	millis := nanos / 1000000
	micros := (nanos - millis*1000000) / 1000
	s += fmt.Sprintf(".%03d,%03d", millis, micros)
	return s
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (f *textFormatter) resolve(line string, e *log.Entry) string {
	// content tags
	line = strings.Replace(line, "${lvlcolor}", lvlColors[e.Level], -1)
	line = strings.Replace(line, "${lvl}", lvlStr(e.Level), -1)
	line = strings.Replace(line, "${timestamp}", timestampStr(e.Time), -1)
	line = strings.Replace(line, "${msg}", e.Message, -1)

	// justify
	if locat := strings.Index(line, "${justify}"); locat >= 0 {
		fill := ""
		if len(e.Data) > 0 {
			tmp := line
			// remove all colors
			for k := range colorCodes {
				tmp = strings.Replace(tmp, k, "", -1)
			}
			if locat = strings.Index(tmp, "${justify}"); locat < justify {
				fill = strings.Repeat(" ", justify-locat)
			}
		}
		line = strings.Replace(line, "${justify}", fill, -1)
	}

	// color tags
	if f.colors {
		for k, v := range colorCodes {
			line = strings.Replace(line, k, v, -1)
		}
	} else {
		for k := range colorCodes {
			line = strings.Replace(line, k, "", -1)
		}
	}

	return line
}
