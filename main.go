package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/anacrolix/torrent"
)

const MeleeMagnet = "magnet:?xt=urn:btih:19c74e4aa2e7306ec739141e662adb5f22b24c64&dn=Super+Smash+Bros.+Melee+%28ENG%29+%5B2014%5D+%7BGameCube+Rom%7D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
const BrawlMagnet = "magnet:?xt=urn:btih:14097be49b11452a86dc2221f8698c7346abdd5f&dn=Super+Smash+Bros+Brawl+%5BUSA%5D+%5BWii%5D+%5BEnglish%5D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
const DolphinPmUrl = "http://download1799.mediafire.com/916h3xkbk2zg/s7q8emzg6ytlttf/Dolphin+5.0-321+SL.zip"
const DolphinFmUrl = "https://www.smashladder.com/downloads/SmashLadder%20Dolphin%20FM-v4.3.zip"

const prompt = `███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗██╗      █████╗ ██████╗ ██████╗ ███████╗██████╗ 
██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║██║     ██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
███████╗██╔████╔██║███████║███████╗███████║██║     ███████║██║  ██║██║  ██║█████╗  ██████╔╝
╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║██║     ██╔══██║██║  ██║██║  ██║██╔══╝  ██╔══██╗
███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║███████╗██║  ██║██████╔╝██████╔╝███████╗██║  ██║
╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝
	0.1-alpha

1) Download Dolphin for Project M
2) Download Dolphin for Melee
3) Download Brawl iso
4) Download Melee iso
e) Exit`

/*
   function formatBytes(bytes,decimals) {
   if(bytes == 0) return '0 Byte';
   var k = 1000;
   var dm = decimals + 1 || 3;
   var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
   var i = Math.floor(Math.log(bytes) / Math.log(k));
   return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}*/

// Convert bytes into human readable sizes
func formatBytes(bytes int) string {
	if bytes == 0 {
		return string(0)
	}
	k := float64(1000)
	b := float64(bytes)
	sizes := [4]string{"bytes", "KB", "MB", "GB"}
	i := math.Floor(math.Log(b) / math.Log(k))
	index := int(i)
	return fmt.Sprintf("%v %v",
		math.Floor((b/math.Pow(k, i))+0.5),
		sizes[index])
}

func DownloadGame(magnet string) error {
	cfg := new(torrent.Config)

	c, err := torrent.NewClient(cfg)
	defer c.Close()
	if err != nil {
		return err
	}
	t, err := c.AddMagnet(MeleeMagnet)
	if err != nil {
		return err
	}
	<-t.GotInfo()
	t.DownloadAll()
	fmt.Println("|", magnet)
	go func() {
		for {
			fmt.Printf("\r| status [%v/%v] ",
				t.BytesMissing(),
				t.BytesCompleted())
		}
	}()
	c.WaitAll()
	log.Print("Done.")
	return err
}

func DownloadDolphin(url, filename string) error {
	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		return err
	}

	log.Printf("Downloading %v", filename)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\r| %v [%v] ", filename, formatBytes(int(n)))
	fmt.Println("Done.")
	return err

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	fmt.Print(">  ")
	switch text, _ := reader.ReadString('\n'); text {
	case "1\n":
		if err := DownloadDolphin(DolphinPmUrl, "dolphin-pm.zip"); err != nil {
			log.Print(err)
		}
	case "2\n":
		if err := DownloadDolphin(DolphinFmUrl, "dolphin-fm.zip"); err != nil {
			log.Print(err)
		}
	case "3\n":
		if err := DownloadGame(BrawlMagnet); err != nil {
			log.Print(err)
		}
	case "4\n":
		if err := DownloadGame(MeleeMagnet); err != nil {
			log.Print(err)
		}
	case "e\n":
		os.Exit(1)
	default:
		log.Print("Invalid input")

	}
}
