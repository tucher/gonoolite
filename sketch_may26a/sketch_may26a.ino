void setup() {
  Serial.begin(9600);
  Serial1.begin(9600);

}

void loop() {
  // put your main code here, to run repeatedly:
  while(Serial.available())
    Serial1.write(Serial.read());
  while(Serial1.available())
    Serial.write(Serial1.read());

}
