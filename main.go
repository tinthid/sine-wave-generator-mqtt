package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	upperBound := 1.0
	lowerBound := -1.0
	noiseFactor := 0.05 // Adjust the noise factor as needed

	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_ADDRESS"))
	opts.SetClientID(os.Getenv("CLIENT_ID"))
	opts.SetUsername(os.Getenv("USERNAME"))
	opts.SetPassword(os.Getenv("PASSWORD"))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to MQTT broker:", token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	start := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		seconds := t.Sub(start).Seconds()
		baseSinValue := ((math.Sin(2*math.Pi*seconds/60)+1)/2)*(upperBound-lowerBound) + lowerBound
		noise := (rand.Float64()*2 - 1) * noiseFactor
		sensorValue := baseSinValue + noise

		if sensorValue > upperBound {
			sensorValue = upperBound
		} else if sensorValue < lowerBound {
			sensorValue = lowerBound
		}

		message := fmt.Sprintf("Time: %s, Sensor Value: %.5f", t.Format("15:04:05"), sensorValue)
		fmt.Println(message)

		token := client.Publish(os.Getenv("TOPIC"), 0, false, message)
		token.Wait()
	}
}
