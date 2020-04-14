package main


import (
	"fmt"
 "log"
  device "github.com/d2r2/go-hd44780"
  "github.com/d2r2/go-i2c"
)

func main() {
	i2c,err:=i2c.NewI2C(0x27,1)
	if err!=nil {log.Fatal(err)}
	lcd,err := device.NewLcd(i2c,device.LCD_16x2)
	if err!=nil {log.Fatal(err)}

	err= lcd.BacklightOn()
	lcd.ShowMessage("aqwerqwer-", device.SHOW_LINE_1)
	fmt.Println("asf")
}
