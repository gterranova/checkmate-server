/*
Copyright © 2023 Gianpaolo Terranova <g.terranova@sazalex.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"terra9.it/checkmate/server/handlers"
	_ "terra9.it/checkmate/server/handlers/auth"
	_ "terra9.it/checkmate/server/handlers/checklist"
	_ "terra9.it/checkmate/server/handlers/html"
	_ "terra9.it/checkmate/server/handlers/json"
	_ "terra9.it/checkmate/server/handlers/markdown"
	_ "terra9.it/checkmate/server/handlers/yaml"

	"terra9.it/checkmate/server/middlewares"
)

var cfgFile string
var hosts map[string]*handlers.Host

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		if _, err := os.Stat(".env"); err == nil {
			viper.SetConfigFile(".env")
			viper.SetConfigType("env")
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("SERVER_HOST", "")
	viper.SetDefault("SERVER_PORT", 4300)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	hosts = make(map[string]*handlers.Host)
	//Get the string that is set in the CONFIG_HOSTS environment variable
	var hostNames = strings.Split(viper.GetString("HOSTS"), " ")
	for _, h := range hostNames {
		hostname := fmt.Sprintf("%s:%d",
			viper.GetString(fmt.Sprintf("HOST_%s_HOSTNAME", h)),
			viper.GetInt("SERVER_PORT"),
		)
		documentFolder := viper.GetString(fmt.Sprintf("HOST_%s_DOCUMENTFOLDER", h))
		hosts[hostname] = &handlers.Host{
			DocumentFolder: strings.ReplaceAll(path.Join(path.Dir("."), documentFolder), "\\", "/"),
		}
	}
}

func New(serverAddr string) {

	//hosts["127.0.0.1:4300"] = &handlers.Host{
	//	DocumentFolder: strings.ReplaceAll(path.Join(path.Dir("."), "..", "documents"), "\\", "/"),
	//}
	//hosts["localhost:4300"] = &handlers.Host{
	//	DocumentFolder: strings.ReplaceAll(path.Join(path.Dir("."), "..", "documents"), "\\", "/"),
	//}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		//EnablePrintRoutes: true,
	})

	// Add CORS Middleware so the frontend get the cookie
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, Content-Type, Session-Id, Authorization",
		MaxAge:           0,
	}))

	key := viper.GetString("AUTH_SECRET")
	app.Use(middlewares.NewAuthMiddleware(key))
	app.Use(middlewares.NewSessionMiddleware(viper.GetString("SERVER_HOST")))

	route := app.Group("/api/v1")

	route.Use(handlers.Page(hosts))
	route.Use(handlers.PageNotFound(hosts))

	app.Static("/", "../www/browser")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile("../www/browser/index.html") })

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	if err := app.Listen(serverAddr); err != nil {
		panic(err)
	}
}

func main() {
	initConfig()
	New(fmt.Sprintf("%v:%v", viper.GetString("SERVER_HOST"), viper.GetInt("SERVER_PORT")))
}
