package main

import (
	xlsx "github.com/tealeg/xlsx/v3"
	"strings"
)

func main() {
	// open an existing file
	wb, err := xlsx.OpenFile("/mnt/c/Users/mengchao.lv1/Desktop/0505.xlsx")
	if err != nil {
		panic(err)
	}
	sh := wb.Sheets[0]

	template := "curl -X POST \"http://qa-ecom-gateway-nlb-1742213771bfcbd8.elb.ap-northeast-1.amazonaws.com/#{url}\" -H \"CHANNEL_ID:user-center-lifecycle\" -H \"accept: application/json\" -H \"Content-Type: application/json\" -d \"#{body}\""

	sh.ForEachRow(func(r *xlsx.Row) error {
		url := r.GetCell(3)
		body := r.GetCell(4)

		curl := r.GetCell(5)

		replace := strings.Replace(template, "/#{url}", url.String(), -1)
		replace = strings.Replace(replace, "#{body}", body.String(), -1)

		curl.SetString(replace)
		return nil
	})

	wb.Save("/mnt/c/Users/mengchao.lv1/Desktop/0505-a.xlsx")
	sh.Close()
}
