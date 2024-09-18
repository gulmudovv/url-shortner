package main

import (
	"fmt"

	"github.com/gulmudovv/url-shortener/internal/config"
)

func init() {

}
func main() {

	cfg := config.MustLoad()

	fmt.Println(cfg)
}
