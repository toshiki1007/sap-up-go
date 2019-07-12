package main

import (
	"fmt"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func calWorkTime(worktimeIn string) string {
	worktimeNum := strings.Split(worktimeIn, "'")
	
	num, _ := strconv.Atoi(worktimeNum[1])

    if int(num)<15{
        worktimeNum[1]="0"
	}else if int(num)<30{
        worktimeNum[1]="25"
	}else if int(num)<45{
        worktimeNum[1]="50"
	}else{
        worktimeNum[1]="75"
	}

    return worktimeNum[0] + "." +worktimeNum[1]
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func main(){
	fmt.Println("******************************************")
	fmt.Println("#事前準備#")
	fmt.Println("1.命名規則に従ったファイルを、実行EXEと同じディレクトリに作成")
	fmt.Println("  命名規則:対象年_対象月_社員番号_センダ原価センタ_勤務WBS(親要素)_勤務WBS(子要素)_休暇WBS(親要素)_休暇WBS(子要素).txt")
	fmt.Println("2.COMPASの勤務時間から日付のある行部分のみをコピーし、1で作成したファイルにペーストする")
	fmt.Println("  ↓ファイルイメージ↓")
	fmt.Println("  1  土  休日  通常      0'00  0'00  0'00   0'00  0'00        小林　謙一  ")
	fmt.Println("  2  日  休日  通常      0'00  0'00  0'00   0'00  0'00        小林　謙一  ")
	fmt.Println("  ・")
	fmt.Println("  ・")
	fmt.Println("  30  日  休日  通常      0'00  0'00  0'00   0'00  0'00        小林　謙一  ")
	fmt.Println("")
	fmt.Println("******************************************")
	
	fmt.Println("ファイル名を入力してください")
	fmt.Println("ex) 2019_06_1210041_26060000_E2600E10_0090_E2826006_0030.txt")
	fmt.Println("")
	fmt.Print("input: ")

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()

	filename := strings.Split(scanner.Text(), ".")[0]

	param := strings.Split(filename, "_")

	if len(param) != 8 {
		fmt.Println("ERROR: 入力ファイル名のフォーマットが不正です")
		os.Exit(1)
	}

	file, err := os.Open("./" + filename + ".txt")
	if err != nil {
		fmt.Println("ERROR: 入力ファイルが存在しません")
		os.Exit(1)
	}
	defer file.Close()

	targetYear := param[0]
	targetMonth := param[1]
	empNo := param[2]
	senderNo := param[3]
	wbsWorkMain := param[4]
	wbsWorkSub := param[5]
	wbsHolidayMain := param[6]
	wbsHolidaySub := param[7]

	holidayList := []string{"Ｓ休", "年休", "Ａ休", "Ｂ休", "結婚", "出産", "忌引", "Ｒ休"}
	holidayHalfList := []string{"前休", "後休"}

	textString := "従業員番号"+"\t"+"センダ原価センタ"+"\t"+"割当先WBS要素"+"\t"+"作業区分"+"\t"+"日付"+"\t"+"作業時間"+"\r\n"
	sjisScanner  := bufio.NewScanner(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	
    for sjisScanner.Scan() {
		line := strings.Split(sjisScanner.Text(), " ")

		if line[4] == "休日"{
			continue
		}

		if line[4] == "振休"{
			textString += empNo +"\t" + senderNo +"\t" + wbsHolidayMain +"\t" + wbsHolidaySub +"\t"
			textString += targetYear + "/" + targetMonth + "/" + line[0] + "\t" + "7.50"     
		}else if contains(holidayList, line[8])  {
			textString += empNo +"\t" + senderNo +"\t" + wbsHolidayMain +"\t" + wbsHolidaySub +"\t"
			textString += targetYear + "/" + targetMonth + "/" + line[0] + "\t" + "7.50"
		}else if contains(holidayHalfList, line[8]) {
			textString += empNo +"\t" + senderNo +"\t" + wbsHolidayMain +"\t" + wbsHolidaySub +"\t"
			textString += targetYear + "/" + targetMonth + "/" + line[0] +"\t" + "3.50"
			textString += "\r\n"
			
			textString += empNo +"\t" + senderNo +"\t" + wbsWorkMain +"\t" + wbsWorkSub +"\t"
			textString += targetYear + "/" + targetMonth + "/" + line[0] +"\t"
			textString += calWorkTime(line[15])
		}else{
			textString += empNo +"\t" + senderNo +"\t" + wbsWorkMain +"\t" + wbsWorkSub +"\t"
			textString += targetYear + "/" + targetMonth + "/" + line[0] +"\t"
			textString += calWorkTime(line[14])
		}

		textString += "\r\n"
	}
	
	file.Close()

    outFile, err := os.OpenFile("./" + filename + "_out.txt", os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
		fmt.Println("ERROR: ファイル出力エラー")
		os.Exit(1)
    }
	defer outFile.Close()
	
    err = outFile.Truncate(0)
    if err != nil {
		fmt.Println("ERROR: ファイル出力エラー")
		os.Exit(1)
	}
	
	sjisWriter  := bufio.NewWriter(transform.NewWriter(outFile, japanese.ShiftJIS.NewEncoder()))
	sjisWriter.WriteString(textString)
	sjisWriter.Flush()

	outFile.Close()

	fmt.Println(filename + "_out.txt を作成しました")
}