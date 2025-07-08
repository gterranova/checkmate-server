package main

//go:generate go run ../tools/mkbundle/main.go

//go:generate fyne bundle --pkg assets --prefix Resource -o assets/assets.go ../settings.chlx
