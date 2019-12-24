package code

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chinfuyang/stock/pkg/path"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

const (
	// TwseEquitiesURL 上市股票代號
	TwseEquitiesURL = "http://isin.twse.com.tw/isin/C_public.jsp?strMode=2"

	// TpexEquitiesURL 上櫃股票代號
	TpexEquitiesURL = "http://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
)

var (
	csvHeader = []string{"type", "code", "name", "ISIN", "start", "market", "group", "CFI"}
)

func fetch(url string) ([][]string, error) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	content := make([][]string, 0)
	var typeStr string
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			content = append(content, csvHeader)
		} else if tr.Children().Length() == 1 {
			typeStr = strings.Trim(big5Conv(tr.Text()), " ")
		} else {
			record := []string{typeStr}
			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				if j == 0 {
					// 將代碼跟名稱拆開
					record = append(record, strings.Split(big5Conv(td.Text()), "　")...)
				} else if j <= 5 {
					// 不抓備註
					record = append(record, big5Conv(td.Text()))
				}
			})
			content = append(content, record)
		}
	})

	return content, nil
}

func big5Conv(needEncoded string) string {
	big5Str, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), needEncoded)

	return big5Str
}

func toCsv(url, path string) error {
	records, err := fetch(url)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
		return err
	}

	return nil
}

// UpdateCodes updates the taiwanese stock code
func UpdateCodes() {
	twsePath := filepath.Join(path.Exec(), "TwseEquities.csv")
	toCsv(TwseEquitiesURL, twsePath)
	tpexPath := filepath.Join(path.Exec(), "TpexEquities.csv")
	toCsv(TpexEquitiesURL, tpexPath)
}
