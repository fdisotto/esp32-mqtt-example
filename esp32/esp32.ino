#include "EspMQTTClient.h"

int LED_BUILTIN = 5;

char WIFI_SSID[] = "";
char WIFI_PASS[] = "";
char MQTT_ADDR[] = "";
char MQTT_USER[] = "";
char MQTT_PASS[] = "";
char MQTT_CLIENT[] = "esp32-client";

char ESP32_STATUS[] = "esp32/status";
char ESP32_LED_STATUS[] = "esp32/led/status";

EspMQTTClient client(
  WIFI_SSID,
  WIFI_PASS,
  MQTT_ADDR,
  MQTT_USER,
  MQTT_PASS,
  MQTT_CLIENT
);

void setup() {
  Serial.begin(115200);
  
  pinMode (LED_BUILTIN, OUTPUT);

  digitalWrite(LED_BUILTIN, HIGH);

  client.enableDebuggingMessages();
  client.enableHTTPWebUpdater();
  client.enableOTA();
  client.enableLastWillMessage(ESP32_STATUS, "off", true);
}

void loop() {
  client.loop();
}

void onConnectionEstablished() {
  client.subscribe(ESP32_LED_STATUS, [] (const String &payload)  {
    if (payload == "on") {
      digitalWrite(LED_BUILTIN, LOW);
    } else {
      digitalWrite(LED_BUILTIN, HIGH);
    }
  });

  client.publish(ESP32_STATUS, "on", true);
}
