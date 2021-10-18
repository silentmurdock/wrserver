package output

import (
	"path"
	"regexp"
	"strconv"
	"strings"
)


type OutputMovieStruct struct {
    Hash string `json:"hash"`
	Quality string `json:"quality"`
	Size string `json:"size"`
	Provider string `json:"provider"`
	Lang string `json:"lang"`
	Title string `json:"title"`
	Seeds string `json:"seeds"`
	Peers string `json:"peers"`
}

type OutputShowStruct struct {
    Hash string `json:"hash"`
	Quality string `json:"quality"`
	Season string `json:"season"`
	Episode string `json:"episode"`
	Size string `json:"size"`
	Provider string `json:"provider"`
	Lang string `json:"lang"`
	Title string `json:"title"`
	Seeds string `json:"seeds"`
	Peers string `json:"peers"`
}

func GetInfoHash(magnet string) string {
	re := regexp.MustCompile(`magnet:\?xt=urn:btih:([a-zA-Z0-9]*)`)
    hash := re.FindAllSubmatch([]byte(magnet), -1)
    if hash == nil {
        return ""
    } else {
    	return string(hash[0][1])
    }
}

func GuessQualityFromString(value string) string {
	// Try to decode quality information from string (url, title, filename)
    lowstr := strings.ToLower(value)	
	quality := ""
	if strings.Contains(lowstr, "1080p") == true {
		quality = "1080p"
	} else if strings.Contains(lowstr, "720p") == true {
		quality = "720p"
	} else if strings.Contains(lowstr, "480p") == true {
		quality = "480p"
	} else if strings.Contains(lowstr, "360p") == true {
		quality = "360p"
	} else if strings.Contains(lowstr, "dvdrip") == true {
		quality = "DVDRIP"
	} else if strings.Contains(lowstr, "bdrip") == true {
		quality = "BDRIP"
	} else if strings.Contains(lowstr, "webrip") == true {
		quality = "WEBRIP"
	} else if strings.Contains(lowstr, "cam") == true {
		quality = "CAM"
	} else {
		quality = "HDTV"
	}
	return quality
}

func DecodeSize(value string) string {
	re := regexp.MustCompile(`[0-9.]+`)
	stringsize := re.FindAllString(value, -1)
	f, _ := strconv.ParseFloat(stringsize[0], 64)
	re = regexp.MustCompile(`(?:GB|MB)`)
	unit := re.FindAllString(value, -1)
	if unit[0] == "GB" {
		f = f * 1024 * 1024 * 1024
	} else if unit[0] == "MB" {
		f = f * 1024 * 1024
	} else if unit[0] == "KB" {
		f = f * 1024
	}
	return strconv.FormatFloat(f, 'f', 0, 64)
}

func DecodeLanguage(value string, language string) string {
	value = strings.TrimSpace(value)
	value = strings.Title(value)
	var enLangArray = [...][2]string{
	    {"ar","Arabic"},{"bg","Bulgarian"},{"hr","Croatian"},{"cs","Czech"},{"da","Danish"},{"nl","Dutch"},{"en","English"},{"et","Estonian"},{"fi","Finnish"},
	    {"fr","French"},{"de","German"},{"el","Greek"},{"he","Hebrew"},{"hu","Hungarian"},{"id","Indonesian"},{"it","Italian"},{"ko","Korean"},{"lv","Latvian"},
	    {"lt","Lithuanian"},{"no","Norwegian"},{"fa","Persian"},{"pl","Polish"},{"pt","Portuguese"},{"ro","Romanian"},{"ru","Russian"},{"sr","Serbian"},{"sk","Slovak"},
	    {"es","Spanish"},{"sw","Swahili"},{"sv","Swedish"},{"th","Thai"},{"tr","Turkish"},{"ur","Urdu"},{"vi","Vietnamese"},
	}

	var huLangArray = [...][2]string{
	    {"ar","Arab"},{"bg","Bolgár"},{"hr","Horvát"},{"cs","Cseh"},{"da","Dán"},{"nl","Holland"},{"en","Angol"},{"et","Észt"},{"fi","Finn"},
	    {"fr","Francia"},{"de","Német"},{"el","Görög"},{"he","Héber"},{"hu","Magyar"},{"id","Indonéz"},{"it","Olasz"},{"ko","Koreai"},{"lv","Lett"},
	    {"lt","Litván"},{"no","Norvég"},{"fa","Perzsa"},{"pl","Lengyel"},{"pt","Portugál"},{"ro","Román"},{"ru","Orosz"},{"sr","Szerb"},{"sk","Szlovák"},
	    {"es","Spanyol"},{"sw","Szuahéli"},{"sv","Svéd"},{"th","Thai"},{"tr","Török"},{"ur","Urdu"},{"vi","Vietnámi"},
	}

	langArray := enLangArray

	switch language {
	case "hu":
		langArray = huLangArray
	}
	
	for _, lang := range langArray {
		if lang[1] == value {
			return lang[0]
		}
	}

	return "en"
}

func RemoveFileExtension(filename string) string {
	return filename[0:len(filename)-len(path.Ext(filename))]
}

func CleanString(value string) string {
	unwanted, err := regexp.Compile(`[^a-zA-Z0-9 _:.+-]+`)
    if err == nil {
        value = unwanted.ReplaceAllString(value, "")
    }

    return strings.TrimSpace(value)
}