package main

import (
	"context"
	spider "github.com/Borislavv/scrapper/internal/spider/app"
	"log"
	"sync"
)

func main() {
	log.Println("starting the spider...")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	//gsh := shutdown.NewGraceful(cancel)

	wg.Add(1)
	go spider.New(ctx).Run(wg)

	//gsh.ListenAndCancel()

	//log.Println("shutting down...")

	wg.Wait()

	cancel()

	log.Println("successfully shut down")
}

//package main

//func main() {
//	defer log.Println("success shutdown!")
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	cfg, err := spiderconfig.Load()
//	if err != nil {
//		panic(err)
//	}
//	taskParser := taskparser.New(cfg)
//	taskProvider := taskprovider.New(cfg, taskParser)
//
//	wg := &sync.WaitGroup{}
//	defer wg.Wait()
//
//	for url := range taskProvider.Provide(ctx) {
//		log.Println("now url: " + url.String())
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//
//			errsCh := make(chan error, 1)
//			msgsCh := make(chan string, 1)
//
//			cwg := &sync.WaitGroup{}
//			defer cwg.Wait()
//
//			cwg.Add(1)
//			go func() {
//				defer cwg.Done()
//				for errr := range errsCh {
//					fmt.Println(errr.Error())
//				}
//			}()
//			cwg.Add(1)
//			go func() {
//				defer cwg.Done()
//				for msg := range msgsCh {
//					fmt.Println(msg)
//				}
//			}()
//
//			defer func() {
//				close(errsCh)
//				close(msgsCh)
//			}()
//
//			innerctx, innercancel := context.WithCancel(ctx)
//			defer innercancel()
//
//			req, rerr := http.NewRequestWithContext(innerctx, http.MethodGet, url.String(), nil)
//			if rerr != nil {
//				panic(rerr)
//			}
//			req.Header.Set("User-Agent", "Chrome/118.0.5993.70 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
//
//			pwg := &sync.WaitGroup{}
//			defer pwg.Wait()
//
//			pwg.Add(1)
//			go doReq(pwg, req, errsCh, msgsCh, 5, innercancel)
//		}()
//	}
//}
//
//func doReq(wg *sync.WaitGroup, req *http.Request, errsCh chan error, msgsCh chan string, attempts int, cancel context.CancelFunc) {
//	defer wg.Done()
//
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		errsCh <- fmt.Errorf("error request: %s: %s", req.URL, err.Error())
//		return
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != 200 {
//		errsCh <- errors.New("bad status code: " + resp.Status)
//
//		if attempts > 0 {
//			errsCh <- errors.New("recursive call")
//			wg.Add(1)
//			go doReq(wg, req, errsCh, msgsCh, attempts-1, cancel)
//		} else {
//			errsCh <- errors.New("all attempts exceeded")
//		}
//
//		return
//	}
//
//	_, err = io.ReadAll(resp.Body)
//	if err != nil {
//		errsCh <- fmt.Errorf("reading all error: " + err.Error())
//		return
//	}
//
//	msgsCh <- fmt.Sprintf("success response: %s", resp.Status)
//
//	cancel()
//}
