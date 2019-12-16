package realtime

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/YangChinFu/stock/pkg/request"
	"github.com/buger/jsonparser"
)

var (
	host = "http://mis.twse.com.tw/"
)

// RealTime 股票現價資訊
type RealTime struct {
	Code                  string   `json:"code"`
	Channel               string   `json:"channel"`
	Name                  string   `json:"name"`
	NameFull              string   `json:"name_full"`
	Time                  string   `json:"time"`
	LatestTradePrice      string   `json:"latest_trade_price"`      // 最近成交價
	TradeVolume           string   `json:"trade_volume"`            // 當盤成交量
	AccumulateTradeVolume string   `json:"accumulate_trade_volume"` // 累積成交量
	Open                  string   `json:"open"`                    // 開盤
	High                  string   `json:"high"`                    // 最高
	Low                   string   `json:"low"`                     // 最低
	BestBidPrice          []string `json:"best_bid_price"`          // 買進五檔價格
	BestBidVolume         []string `json:"best_bid_volume"`         // 買進五檔數量
	BestAskPrice          []string `json:"best_ask_price"`          // 賣出五檔價格
	BestAskVolume         []string `json:"best_ask_volume"`         // 賣出五檔數量
}

func buildURL(path string, query string) string {
	return fmt.Sprintf("%s%s?_=%d&ex_ch=%s", host, path, time.Now().Unix(), query)
}

// Get 回傳搜尋的股票代碼即時資訊
func Get(codes ...string) ([]RealTime, error) {
	url := buildURL("stock/api/getStockInfo.jsp", combineCode(codes))
	resp, err := request.Get(url)

	if err != nil {
		return nil, err
	}

	realtimes, err := formatStockData(resp)

	if err != nil {
		return nil, err
	}

	return realtimes, nil
}

func combineCode(codes []string) string {
	var tmpStrs []string
	for _, e := range codes {
		tmpStrs = append(tmpStrs, "tse_"+e+".tw")
	}
	sort.Strings(tmpStrs)

	return strings.Join(tmpStrs, "|")
}

func formatStockData(data []byte) ([]RealTime, error) {
	content, _, _, err := jsonparser.Get(data, "msgArray")
	if err != nil {
		log.Fatal(err)

		return nil, err
	}
	var stocks []RealTime

	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		stocks = append(stocks, newRealTime(value))
	})

	return stocks, nil
}

func splitBest(str string, err error) []string {
	return strings.Split(str, "_")
}

func newRealTime(data []byte) RealTime {
	realTime := RealTime{}
	realTime.Code, _ = jsonparser.GetString(data, "c")
	realTime.Channel, _ = jsonparser.GetString(data, "ch")
	realTime.Name, _ = jsonparser.GetString(data, "n")
	realTime.NameFull, _ = jsonparser.GetString(data, "nf")
	strTimestamp, _ := jsonparser.GetString(data, "tlong")
	timestamp, _ := strconv.ParseInt(strTimestamp, 10, 64)
	tm := time.Unix(timestamp/1000, 0)
	realTime.Time = tm.Format("2006-01-02 03:04:05 PM")
	realTime.AccumulateTradeVolume, _ = jsonparser.GetString(data, "v")
	realTime.LatestTradePrice, _ = jsonparser.GetString(data, "z")
	realTime.TradeVolume, _ = jsonparser.GetString(data, "tv")
	realTime.Open, _ = jsonparser.GetString(data, "o")
	realTime.High, _ = jsonparser.GetString(data, "h")
	realTime.Low, _ = jsonparser.GetString(data, "l")
	realTime.BestBidPrice = splitBest(jsonparser.GetString(data, "b"))
	realTime.BestBidVolume = splitBest(jsonparser.GetString(data, "g"))
	realTime.BestAskPrice = splitBest(jsonparser.GetString(data, "a"))
	realTime.BestAskVolume = splitBest(jsonparser.GetString(data, "f"))

	return realTime
}
