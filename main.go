package main

import (
    "github.com/ChimeraCoder/anaconda"
    "math/rand"
    "time"
    "fmt"
    "log"
)

func main() {
    fmt.Println("Program initiated...")
    api := createApi()
    remove(api)
    follow(api, 822011433470738432, 150)
}

/*
follow automatically
 */
func follow(api *anaconda.TwitterApi, id int64, max int) {
    fmt.Println("follow start")
    users, _ := api.GetFollowersUser(id, nil)
    rand.Seed(time.Now().UnixNano())
    for i:=0; i<max; i ++ {
        randomId := users.Ids[rand.Intn(len(users.Ids))]

        _, err := api.FollowUserId(randomId, nil)
        if err != nil {
            log.Println(err)
        }
        time.Sleep(1 * time.Second)
    }
    fmt.Println("follow finished")
}

func remove(api *anaconda.TwitterApi) {
    f, err := api.GetFriendsIds(nil)
    if err != nil {
        log.Println(err)
    }

    fed, err := api.GetFollowersIds(nil)
    if err != nil {
        log.Println(err)
    }

    result := filter(f.Ids, fed.Ids)
    for _, i := range result {
        api.UnfollowUserId(i)
    }
    time.Sleep(1000 * time.Millisecond)
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

func createApi() *anaconda.TwitterApi {
    anaconda.SetConsumerKey("YOUR_CONSUMER_KEY")
    anaconda.SetConsumerSecret("YOUR_CONSUMER_KEY")
    api := anaconda.NewTwitterApi("YOUR_ACCESS_TOKEN", "YOUR_ACCESS_SECRET")
    return api
}