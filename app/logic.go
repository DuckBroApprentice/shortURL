package app

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/redis/go-redis/v9"
)

var Ctx context.Context

func RandomNum() int {
	base := rand.Intn(916132832) //62^5
	return base
}

// a-z,A-Z,0-9
func Base62(randomNum int) string {
	trans := make(map[int]string)
	trans = map[int]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "a",
		11: "b",
		12: "c",
		13: "d",
		14: "e",
		15: "f",
		16: "g",
		17: "h",
		18: "i",
		19: "j",
		20: "k",
		21: "l",
		22: "m",
		23: "n",
		24: "o",
		25: "p",
		26: "q",
		27: "r",
		28: "s",
		29: "t",
		30: "u",
		31: "v",
		32: "w",
		33: "x",
		34: "y",
		35: "z",
		37: "A",
		38: "B",
		39: "C",
		40: "D",
		41: "E",
		42: "F",
		43: "G",
		44: "H",
		45: "I",
		46: "J",
		47: "K",
		48: "L",
		49: "M",
		50: "N",
		51: "O",
		52: "P",
		53: "Q",
		54: "R",
		55: "S",
		56: "T",
		57: "U",
		58: "V",
		59: "W",
		60: "X",
		61: "Y",
		62: "Z",
	}
	if randomNum/62 == 0 {
		num := randomNum % 62
		if num == 0 {
			return ""
		}
		return trans[num]
	}

	return Base62(randomNum/62) + trans[randomNum%62]
}

func InputURL(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm() err: %v", err)
	}
	//參考別人的作法
	originalURL, err := url.ParseRequestURI(r.Form.Get("url"))
	if err != nil || originalURL.Scheme == "" || originalURL.Host == "" {
		//panic(err)
	}
	outputInt := RandomNum()
	outputURL := Base62(RandomNum())
	fmt.Println(outputURL)
	//check redis
	val, err := client.Get(Ctx, outputURL).Result()
	fmt.Println("After Get outputURL is :", outputURL)      //correct 但會發生同一個網址不同結果 要修正random的起點
	fmt.Println("After Get outputURL become val is :", val) //沒有set過的key對應的value一定是nil  correct
	fmt.Println("After Get err is :", err)                  //correct
	if err != nil {
		log.Println("val == nil prepare set new key")
		if err == redis.Nil {
			set_err := client.Set(Ctx, outputURL, originalURL.String(), 0).Err() //這裡會invalid memory address or nil pointer dereference
			if set_err != nil {
				panic(set_err)
			}
		}
	}
	if err == nil && val != originalURL.RawQuery {
		outputURL = Base62(outputInt * 62)
	}

	fmt.Printf("localhost:8080/short?short=%s", outputURL)
}

func Redirect(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	query := r.URL.Query()
	shortURL := query.Get("short")
	log.Printf("Func Redirect Get shortURL: %s\n", shortURL)
	longURL, err := client.Get(Ctx, shortURL).Result()
	log.Printf("success Get longURL: %s", longURL)
	if err == redis.Nil {
		http.Redirect(w, r, "/index", 404)
	} else if err != nil {
		log.Println("資料庫錯誤")
	}
	log.Printf("shortURL is exist key in redis")
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
