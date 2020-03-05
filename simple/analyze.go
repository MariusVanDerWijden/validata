package simple

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/MariusVanDerWijden/mvp/validata"
	"github.com/montanaflynn/stats"
)

func AnalyzeNum(column int, data []float64) error {
	log.Print("Analyzing: ", data)
	outliers, err := stats.QuartileOutliers(data)
	if err != nil {
		return err
	}
	if len(outliers.Extreme) != 0 {
		fmt.Println("found outlier")
		msgs := make([]string, len(outliers.Extreme))
		for i, ex := range outliers.Extreme {
			var row int
			for i, ele := range data {
				if ele == ex {
					row = i
				}
			}
			msgs[i] = fmt.Sprintf("Field in column %d, row %d with value %f is an extreme outlier", column, row, ex)
		}
		Notify(msgs)
		os.Exit(0)
	} else if len(outliers.Mild) != 0 {
		fmt.Println("found mild outliers")
	} else {
		fmt.Println("no outliers found")
	}
	return nil
}

func Notify(msgs validata.Messages) {
	mailBody := "Dear Sir or Madam,\n\nin the previous 10 seconds, the following Datasets have been flagged by our automated system, because:\n"
	for _, msg := range msgs {
		mailBody += "\t" + msg + "\n"
	}
	mailBody += "You can check these issues by following this link www.vali.data/asdf\n\nBest Regards\nThe vali.data team"
	sendEmail(mailBody)
}

func sendEmail(body string) {
	log.Println("Sending email")
	auth := smtp.PlainAuth(
		"",
		"bot.vali.data@gmail.com",
		"bushdid911",
		"smtp.gmail.com",
	)
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"bot.vali.data@gmail.com",
		[]string{"steffan-alex@web.de", "marius@sm2.network"},
		[]byte("To: marius@sm2.network\r\n"+
			"Subject: Alert: vali.data found potential invalid data\r\n"+
			"\r\n"+
			body),
	)
	if err != nil {
		log.Fatal(err)
	}
}
