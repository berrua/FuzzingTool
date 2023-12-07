/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var fuzzerCmd = &cobra.Command{
	Use:   "fuzzer",
	Short: "Web Fuzzer Tool",
	Long:  `Welcome to the Web Fuzzer Tool! You can use the '--url' and '--wordlist' parameters to fuzz any targeted website.`,
	Run: func(cmd *cobra.Command, args []string) {
		//kullanıcıdan parametreler alınır, değişkenlere atanır
		fuzzUrl, _ := cmd.Flags().GetString("url")
		fuzzWordlist, _ := cmd.Flags().GetString("wordlist")
		fuzzSpeed, _ := cmd.Flags().GetInt("speed")
		fuzzStatus, _ := cmd.Flags().GetString("status")
		fuzzTimeout, _ := cmd.Flags().GetInt("timeout")

		if fuzzUrl == "" || fuzzWordlist == "" {
			c := color.New(color.FgMagenta, color.Bold)
			c.Println("Please enter a URL using the format '--url=<URL>' and provide a wordlist path using '--wordlist=wordlist.txt'")
		} else {
			c := color.New(color.FgGreen, color.Bold)
			c.Println("Scanning...")

			file, err := os.Open(fuzzWordlist)
			if err != nil {
				c := color.New(color.FgRed, color.Bold)
				c.Println("Error opening wordlist file:", err)
				os.Exit(1)
			}
			scanner := bufio.NewScanner(file) //scanner oluşturulur

			var wg sync.WaitGroup          //waitgroup oluşturulur, workerların tamamlanmasını beklemek için kullanılır
			jobs := make(chan string, 100) //workerlara gönderilecek işlerin(kelimelerin) kanalı

			//workerlar başlatılır ve her worker, worker fonksiyonunu çağırır
			for i := 0; i < 3; i++ {
				wg.Add(1)
				go worker(fuzzUrl, fuzzStatus, fuzzTimeout, fuzzSpeed, &wg, jobs)
			}

			//scanner ile wordlist okunur ve işler(kelimeler) workerlara gönderilir
			for scanner.Scan() {
				jobs <- scanner.Text()
			}

			close(jobs) //kanal kapatılır, tüm işlerin(kelimelerin) işlendiğini workerlara bildirir

			wg.Wait() //tüm workerların tamamlanması beklenir

			c = color.New(color.FgGreen, color.Bold)
			c.Println("Completed!")
		}
	},
}

func worker(fuzzUrl string, fuzzStatus string, fuzzTimeout int, fuzzSpeed int, wg *sync.WaitGroup, jobs <-chan string) {
	defer wg.Done() //fonksiyon tamamlandığında waitgroupa bildirir

	//HTTP client ile zaman aşımı süresi ayarlanır
	client := http.Client{
		Timeout: time.Duration(fuzzTimeout) * time.Millisecond,
	}

	for job := range jobs {
		//hedef URLe iş(kelime) eklenir ve istek gönderilir
		resp, err := client.Get(fuzzUrl + "/" + job)
		//istek başarısız olursa hata mesajı verilir ve program sonlanır
		if err != nil {
			c := color.New(color.FgRed, color.Bold)
			c.Println("Error making HTTP request:", err)
			os.Exit(1)
		}

		//durum kodu kontrol edilir
		if resp.StatusCode != 404 {
			statusCodeStr := fmt.Sprint(resp.StatusCode)     //durum kodu değişkene atanır
			if strings.Contains(statusCodeStr, fuzzStatus) { //status parametresinin kullanımı kontrol edilir
				color.New(color.FgMagenta, color.Bold).Print(job)
				fmt.Print(" - ")
				fmt.Print(fuzzUrl+"/"+job, " - ")
				color.New(color.FgBlue, color.Bold).Println(resp.StatusCode)
			}
		}
		resp.Body.Close()

		time.Sleep(time.Duration(fuzzSpeed) * time.Millisecond) //verilen zaman miktarına kadar çalışmayı askıya alır
	}
}

func init() {
	displayAsciiArt()

	rootCmd.AddCommand(fuzzerCmd)

	//fuzzerCmd komutunun parametreleri
	fuzzerCmd.PersistentFlags().String("url", "", "Specify the target URL for fuzzing")                              //hedeflenen URL
	fuzzerCmd.PersistentFlags().String("wordlist", "", "Specify the name of the wordlist file")                      //wordlist tercihi
	fuzzerCmd.PersistentFlags().String("status", "", "Specify the status code to filter the results (ex: 200, 403)") //HTTP durum kodları filtrelenir
	fuzzerCmd.PersistentFlags().Int("speed", 500, "Specify the speed in milliseconds for fuzzing")                   //fuzzing işlemi her 500 milisaniyelik bir hızda çalışır
	fuzzerCmd.PersistentFlags().Int("timeout", 5000, "Specify the timeout in milliseconds for HTTP requests")        //HTTP isteklerinin zaman aşımı süresini belirtir
}

func displayAsciiArt() {
	file, err := os.ReadFile("ascii_art.txt")
	if err != nil {
		c := color.New(color.FgRed, color.Bold)
		c.Println("ASCII art could not be displayed:", err)
		return
	}
	fmt.Println(string(file))
}
