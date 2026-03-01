TARGET = weather
SRC = ./src
BUILD = ./build
OUT = ./out

all: build runall

build:
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(TARGET) $(SRC)

runall: build
	mkdir -p $(OUT)
	$(BUILD)/$(TARGET)  35.9653 -83.9233 | jgraph -P | ps2pdf -sPAPERSIZE=letter -dNOEPS -dEPSCrop - > $(OUT)/knoxville.pdf
	$(BUILD)/$(TARGET)  51.5072 -0.12760 | jgraph -P | ps2pdf -sPAPERSIZE=letter -dNOEPS -dEPSCrop - > $(OUT)/london.pdf
	$(BUILD)/$(TARGET)  64.8401 -147.720 | jgraph -P | ps2pdf -sPAPERSIZE=letter -dNOEPS -dEPSCrop - > $(OUT)/fairbanks.pdf
	$(BUILD)/$(TARGET) -31.9514 115.8617 | jgraph -P | ps2pdf -sPAPERSIZE=letter -dNOEPS -dEPSCrop - > $(OUT)/perth.pdf
	$(BUILD)/$(TARGET)  17.4065  78.4772 | jgraph -P | ps2pdf -sPAPERSIZE=letter -dNOEPS -dEPSCrop - > $(OUT)/hyderabad.pdf

clean:
	rm -rf $(BUILD) $(OUT)
