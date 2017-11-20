package main

import (
    "github.com/ChimeraCoder/anaconda"
    "math/rand"
    "time"
    "net/url"
    "strconv"
    "github.com/robfig/cron"
)

var stopper = make(chan bool)

func main() {
    //var id int64 = 4001665694
    api := createApi()
    c := cron.New()
    c.AddFunc("@daily", func() {follow(api, 884680626334334976, 150)})
    c.AddFunc("0 0 12 * * *", func() {c.AddFunc("@daily", func() {remove(api)})})
    c.Start()
    <- stopper
}

/*
follow automatically
 */
func follow(api *anaconda.TwitterApi, id int64, max int) {
    users, _ := api.GetFollowersUser(id, nil)
    rand.Seed(time.Now().UnixNano())
    for i:=0; i<max; i ++ {
        randomId := users.Ids[rand.Intn(len(users.Ids))]
        _, err := api.FollowUserId(randomId, nil)
        if err != nil {
            i -= 1
        }
    }
}

func remove(api *anaconda.TwitterApi) {
    f, _ := api.GetFriendsIds(nil)
    for _, id :=  range f.Ids {
        if !isFollowMe(api, id) {
            api.UnfollowUserId(id)
        }
    }
}

func isFollowMe(api *anaconda.TwitterApi, userId int64) bool {
    v := url.Values{"user_id": {strconv.FormatInt(userId, 10)}}
    user, _ := api.GetFriendsIds(v)
    for _, id := range user.Ids {
        if id == 932137292197593089 {
            return true
        }
    }
    return false
}

func createApi() *anaconda.TwitterApi {
    anaconda.SetConsumerKey("3rJOl1ODzm9yZy63FACdg")
    anaconda.SetConsumerSecret("5jPoQ5kQvMJFDYRNE8bQ4rHuds4xJqhvgNJM4awaE8")
    api := anaconda.NewTwitterApi("932137292197593089-kylZz2L7TmKdn6Tdzuf6yH9uFEs7s2Y", "wcmU6LH10C1sEChz9X3jexjDaP5UIoTKfRqQVZpGmwFye")
    return api
}