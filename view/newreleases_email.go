package view

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"

	"github.com/ynori7/music/allmusic"
)

type HtmlTemplate struct {
	Discographies []allmusic.Discography
}

func NewHtmlTemplate(discographies []allmusic.Discography) HtmlTemplate {
	return HtmlTemplate{
		Discographies: discographies,
	}
}

func (h HtmlTemplate) ExecuteHtmlTemplate() (string, error) {
	t := template.Must(template.New("html").
		Funcs(template.FuncMap{
			"mod": func(i, j int) bool { return i%j == 0 },
			"allmusicRating": func(r int) string {
				switch r {
				case 0, 1:
					return "https://cdn-gce.allmusic.com/images/global/rating/allmusic-0.png"
				default:
					return fmt.Sprintf("https://cdn-gce.allmusic.com/images/global/rating/allmusic-%d.png", r-1)
				}
			},
			"coverImage": func(s string) string {
				if len(s) == 0 {
					return "https://cdn-gce.allmusic.com/images/no_image/album_300x300.png"
				}
				return s
			},
		}).
		Parse(htmlTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := t.Execute(w, h)
	if err != nil {
		return "", err
	}

	w.Flush()
	return b.String(), nil
}

const htmlTemplate = `<html>
<head>
	<style>
	a {
		text-decoration: none;
	}
	h3.artist {
		font-size:10pt;line-height:11pt;width:100%;margin:5px 0;padding-top:5px;font-weight:normal;
	}
	h3.release {
		font-size:11pt;margin:5px 0px;padding:0px;
	}
	h4.genres {
		font-size:9pt;margin:7px 0;padding:0;
	}
	h4.genres span {
		text-decoration:none;color:#000;font-weight:normal;width:100%;
	}
	p.score {
		font-size:10pt;line-height:14.5pt;margin:10px 0 20px;
	}
	</style>
</head>
<body>
<table width="640" cellpadding="0" cellspacing="0" border="0" bgcolor="#FFFFFF">
    <tbody><tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
    <tr>
        <td align="center" valign="top">
            <table width="600" cellpadding="0" cellspacing="0" border="0" style="border-bottom:1px;border-bottom-style:solid;border-color:#e3e3e3">
                <tbody>
				{{range $i, $val := .Discographies}}
					{{ if eq $i 0 }}<tr>{{ else if mod $i 3 }}</tr><tr>{{ else }}<td width="2%" align="center" valign="top">&nbsp;</td>{{ end }}
					<td width="32%" align="left" valign="top">
                    	<a href="{{ $val.NewestRelease.Link }}">
                        	<img style="width:100%" src="{{ coverImage $val.NewestRelease.Image }}"><br>
                        </a>
                        <h3 class="artist">
                        	<a href="{{ $val.Artist.Link }}">{{ $val.Artist.Name }}</a>
                        </h3>
                        <h3 class="release">
                        	<a href="{{ $val.NewestRelease.Link }}">{{ $val.NewestRelease.Title }}</a>
						</h3>
                        <h4 class="genres">
							<span>
							{{ range $j, $genre := $val.Artist.Genres }}{{ if ne $j 0 }}, {{ end }}{{ $genre }}{{ end }}
							</span>
                        </h4>
						<img src="{{ allmusicRating $val.NewestRelease.Rating }}" width="auto" height="auto" alt="star rating"><br>
                        <p class="score">Score: {{ $val.Score }} / 10</p>
                        </td>
				{{ end }}
              </tr>
            </tbody></table>
        </td>
    </tr>
    <tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
</tbody></table>
</body>
</html>
`
