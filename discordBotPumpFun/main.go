package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/chromedp"
)

const (
	pumpFunURL   = "https://pump.fun/board"
	discordToken = "" // Replace with your bot token
	channelID    = ""
)

var (
	seenCoins = make(map[string]bool)
	mutex     sync.Mutex
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting bot...")

	sess, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	if err := sess.Open(); err != nil {
		log.Fatalf("Error connecting to Discord: %v", err)
	}
	defer sess.Close()

	log.Println("Bot successfully connected to Discord")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Основной цикл
	for {
		log.Println("Starting new check cycle...")
		if err := checkNewCoins(ctx, sess); err != nil {
			log.Printf("Error during check: %v", err)
		}
		log.Println("Check cycle completed, waiting for next iteration...")
		time.Sleep(5 * time.Second)
	}
}

func checkNewCoins(ctx context.Context, sess *discordgo.Session) error {
	var coins []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(pumpFunURL),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
            Array.from(document.querySelectorAll('div[data-sentry-component="CoinPreview"]')).map(el => {
                const title = el.querySelector('p.text-sm span.font-bold')?.textContent || '';
                const marketCap = el.querySelector('div.text-xs.text-green-300')?.textContent || '';
                const timeAgo = el.querySelector('span.w-full.xl\\:w-auto')?.textContent || '';
                const replies = el.querySelector('p:-soup-contains("replies")')?.textContent || '';
                const creator = el.querySelector('button span')?.textContent || '';
                const id = el.id || '';
                return JSON.stringify({title, marketCap, timeAgo, replies, creator, id});
            })
        `, &coins),
	)

	if err != nil {
		return fmt.Errorf("failed to execute Chrome actions: %v", err)
	}

	log.Printf("Found %d coins", len(coins))

	for _, coinJSON := range coins {
		var coin struct {
			Title     string `json:"title"`
			MarketCap string `json:"marketCap"`
			TimeAgo   string `json:"timeAgo"`
			Replies   string `json:"replies"`
			Creator   string `json:"creator"`
			ID        string `json:"id"`
		}

		if err := json.Unmarshal([]byte(coinJSON), &coin); err != nil {
			log.Printf("Error parsing coin data: %v", err)
			continue
		}

		title := strings.TrimSpace(strings.Split(coin.Title, ":")[0])
		if title == "" {
			continue
		}

		mutex.Lock()
		if seenCoins[title] {
			mutex.Unlock()
			continue
		}
		seenCoins[title] = true
		mutex.Unlock()

		log.Printf("New coin found: %s", title)

		message := fmt.Sprintf("🚀 New Coin Alert! 🚀\n"+
			"**%s**\n"+
			"👨‍💼 Created by: %s\n"+
			"💰 %s\n"+
			"⏰ %s\n"+
			"💬 %s\n"+
			"🔗 https://pump.fun/coin/%s",
			title, coin.Creator, coin.MarketCap, coin.TimeAgo, coin.Replies, coin.ID)

		if _, err := sess.ChannelMessageSend(channelID, message); err != nil {
			log.Printf("Error sending message to Discord: %v", err)
		} else {
			log.Printf("Successfully sent message for coin: %s", title)
		}
	}

	return nil
}
