package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/v2fly/vmessping"
	"context"

	"github.com/xxf098/lite-proxy/web"
)

func main() {
	v2rayConfigFiles := "https://raw.githubusercontent.com/barry-far/V2ray-Configs/main/All_Configs_Sub.txt"

	configFileName := "configs.txt"

	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("file exists")
		errFile := os.Remove(configFileName)

		if errFile != nil {
			log.Fatal(errFile)
			os.Exit(0)
		}
	}

	out, err := os.Create(configFileName)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	defer out.Close()

	fmt.Println("create file")

	resp, err := http.Get(v2rayConfigFiles)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-2)
	}
	defer resp.Body.Close()

	fmt.Print("download from web:")
	fmt.Println(resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	_, err1 := io.Copy(out, resp.Body)

	if err1 != nil {
		log.Fatal(err.Error())
		os.Exit(-2)
	}

	fmt.Println("download file and copy is complete")

	file, err := os.OpenFile(configFileName, os.O_RDONLY, os.ModePerm)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-3)
	}

	// const maxCapacity int = 200
	// buf := make([]byte, maxCapacity)

	scanner := bufio.NewScanner(file)

	// scanner.Buffer(buf, maxCapacity)
	fmt.Println("start scan")

	outputConfigs := make([]string, 0)

	for scanner.Scan() {

		line := scanner.Text()

		_, err := url.ParseRequestURI(line)

		if err != nil {
			continue
		}

		ctx := context.Background()

		opts := web.ProfileTestOptions{
			GroupName:     "Default",
			SpeedTestMode: "pingonly",   //  pingonly speedonly all
			PingMethod:    "googleping", // googleping
			SortMethod:    "rspeed",     // speed rspeed ping rping
			Concurrency:   300,
			TestMode:      2,
			Subscription:  line,
			Language:      "en", // en cn
			FontSize:      24,
			Theme:         "rainbow",
			Unique:        true,
			Timeout:       5 * time.Second,
			OutputMode:    0,
		}

		nodes, errLink := web.TestContext(ctx, opts, &web.EmptyMessageWriter{})
		if errLink != nil {
			fmt.Println("error in this line")
			fmt.Println(errLink)
		}

		for _, node := range nodes {
			ping, errPing := strconv.ParseInt(node.Ping, 10, 32)
			if errPing != nil {
				continue
			}

			if node.IsOk && ping > 0 {
				outputConfigs = append(outputConfigs, node.Link)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		os.Exit(-4)
	}

	f, errFile := os.Create("./valid.txt")

	if errFile != nil {
		fmt.Println("error in write file")
		fmt.Println(errFile)
	}

	f.WriteString(strings.Join(outputConfigs[:], "\n"))
	f.Sync()

	defer f.Close()

}
