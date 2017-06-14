# build dictionary

install: ipa_internal uni_internal bindata ipa uni
	-rm ../../internal/dic/data/bindata*.go
	cp ./bindata*.go ../../internal/dic/data
	cp ipa.dic ../../_sample
	cp uni.dic ../../_sample

.PHONY: ipa_internal
ipa_internal:
	-mkdir -p ./dic/ipa
	go run main.go ipa -mecab ./mecab-ipadic-2.7.0-20070801 -output ./dic/ipa/ -z=false

.PHONY: ipa
ipa:
	go run main.go ipa -mecab ./mecab-ipadic-2.7.0-20070801

.PHONY: uni_internal
uni_internal:
	-mkdir -p ./dic/uni
	go run main.go uni -mecab ./unidic-mecab-2.1.2_src -output ./dic/uni/ -z=false

.PHONY: uni
uni:
	go run main.go uni -mecab ./unidic-mecab-2.1.2_src

.PHONY: bindata
bindata: 
	go-bindata -o bindata.go -nomemcopy -separate -pkg=data dic/...

.PHONY: clean
clean:
	-rm ./bindata*.go
	-rm ipa.dic
	-rm -rf ./dic

