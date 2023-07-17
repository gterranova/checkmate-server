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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"

	"terra9.it/vadovia/assets"
	"terra9.it/vadovia/core"
)

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".env")
	}

	viper.SetDefault("API_SERVER_HOST", "")
	viper.SetDefault("API_SERVER_PORT", 4300)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func ProjectData(ctx *fiber.Ctx) error {
	project := core.NewProject()
	return ctx.JSON(
		fiber.Map{
			"project": project,
			"result":  project.Describe(),
		},
	)
}

func SetupRoutes(app *fiber.App) {

	app.Get("/assets/logo.png", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(assets.ResourceLogovadoviacompactPng.StaticContent)
	})

	route := app.Group("/api/v1")

	route.Get("/project", ProjectData)
}

func New(serverAddr string) {

	app := fiber.New()

	// Add CORS Middleware so the frontend get the cookie
	app.Use(cors.New())

	SetupRoutes(app)

	app.Listen(serverAddr)
}

func main() {
	initConfig()
	New(fmt.Sprintf("%v:%v", viper.GetString("API_SERVER_HOST"), viper.GetInt("API_SERVER_PORT")))
}
