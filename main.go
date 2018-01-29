package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/skratchdot/open-golang/open"
)

var (
	api            *anaconda.TwitterApi
	ConsumerKey          = "CK"
	ConsumerSecret       = "CS"
	ID             int64 = 878237273061326848
	selfID         int64
)

func init() {
	// Twitter 側の問題で golang では http/2 を無効にしないとエラーメッセージが文字化けする。
	godebug := os.Getenv("GODEBUG")
	if godebug != "" {
		godebug += ","
	}
	godebug += "http2client=0"
	os.Setenv("GODEBUG", godebug)
}

func main() {
	a, err := getAuth()
	api = a
	if err != nil {
		log.Fatal(err)
	}
	selfData, _ := api.GetSelf(nil)
	selfID = selfData.Id
	t := getRemoveTarget()
	fmt.Printf("[x]片思いの数は%v人です\n", len(t))
	//remove(t)
	follow(878237273061326848, 500)
}

/*
follow automatically
*/
func follow(id int64, max int) {
	users, _ := api.GetFollowersUser(id, nil)
	var tha = []int64{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < max; i++ {
		randomId := users.Ids[rand.Intn(len(users.Ids))]
		if cck(randomId, tha) {
			tha = append(tha, randomId)
		} else {
			i--
		}
	}

	for n, i := range tha {
		_, err := api.FollowUserId(i, nil)
		fmt.Printf("[+]%v人目のフォローをしました\n", n+1)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("follow finished")
}

func cck(r int64, tha []int64) bool {
	for _, t := range tha {
		if t == r || r == selfID {
			return false
		}
	}
	return true
}

type Error struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func remove(t []int64) {
	for r, i := range t {
		_, err := api.UnfollowUserId(i)
		fmt.Printf("[-]%v人目をリムーブしました\n", r+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func getRemoveTarget() []int64 {
	f, err := api.GetFriendsIds(nil)
	if err != nil {
		log.Println(err)
	}

	fed, err := api.GetFollowersIds(nil)
	if err != nil {
		log.Println(err)
	}

	return filter(f.Ids, fed.Ids)
}

func filter(lhs, rhs []int64) []int64 {
	m := map[int64]int{}

	for _, v := range lhs {
		if _, ok := m[v]; !ok {
			m[v] = 1
		}
	}

	for _, v := range rhs {
		if _, ok := m[v]; ok {
			m[v] = 2
		}
	}

	var ret []int64

	for i := range m {
		if m[i] == 1 {
			ret = append(ret, i)
		}
	}
	return ret
}

func getAuth() (*anaconda.TwitterApi, error) {
	// ランダムなポートで受け待つ
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	defer l.Close()

	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)

	// authorizationURL にアクセスして認証 URI を得る
	uri, cred, err := anaconda.AuthorizationURL(fmt.Sprintf("http://%s/", l.Addr().String()))
	if err != nil {
		return nil, err
	}
	open.StartWith(uri, "safari")

	// コールバックを用意する
	var verifier string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		verifier = r.URL.Query().Get("oauth_verifier")
		// リスナを閉じて処理を続行
		l.Close()
	})

	http.Serve(l, nil)

	// 認証
	cred, _, err = anaconda.GetCredentials(cred, verifier)
	if err != nil {
		return nil, err
	}

	// API クライアントを設定
	api := anaconda.NewTwitterApi(cred.Token, cred.Secret)
	fmt.Println("authorize successful")
	return api, nil
}
