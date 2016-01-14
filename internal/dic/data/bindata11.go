package data

import(
	"os"
	"time"
)

var _dicUniUnkDic = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xe2\xfd\xdf\xc9\xc4\xc8\xf4\xbf\x8b\x81\xf1\x7f\x07\x03\x83\xe9\xff\x76\x66\x46\x46\x56\xdf\xfc\xa2\x82\x0c\x90\x00\x23\x33\x23\x9b\x4f\x6a\x5a\x89\xa7\x0b\x23\x0b\x03\x23\x7b\x50\x66\x7a\x06\x94\xcd\x16\x9e\x0a\xe2\x00\x99\x0c\x0c\xff\x18\x4f\x03\x0d\x50\x66\xfc\xa7\x9a\x04\xc1\x59\x19\x0c\x8c\xff\x34\xd4\x20\x38\xea\x11\x90\xa3\x93\x07\xc1\xe5\x32\x40\x8e\x6a\x1e\x04\x57\x27\x81\x38\x69\x10\x9c\xaf\x81\xc4\xc9\xcc\x63\x40\x98\x96\xf6\x06\xc9\xb4\x70\x0e\x24\x03\x8a\x6e\x00\x39\x7a\x0b\x20\x58\x4e\x09\x49\xa6\x44\x06\x49\x4f\x58\x18\x92\x69\x19\x41\x48\xf6\x64\x9d\x03\x29\x13\x82\xe0\x39\x11\x48\xca\xfa\xae\x20\x99\x36\xff\x0b\x92\x17\xba\xbe\x20\x19\xd0\xb7\x0b\xc9\x9e\x5a\x09\x24\x4e\x51\x13\x92\xd1\xed\x16\x48\xa6\x4d\xfd\x83\x6c\x0f\x17\x92\x69\xbd\x3e\x20\x4e\x09\x04\xe7\x6d\x41\x52\x96\x29\x86\x64\x74\xe8\x09\x24\xa3\x8b\x97\x20\x19\x90\x5d\x87\x64\x4f\xe1\x16\xe4\x80\x0f\x40\x32\x3a\xd1\x05\x14\x6e\x73\x20\x58\x3f\x00\xc9\xb4\xd6\x1e\x06\xbe\xff\x8d\x2c\xc0\x14\x01\x74\x3e\x0b\x28\x76\xa5\x80\x2c\x6e\x11\x2e\x06\x21\x01\x23\x0e\x07\x36\x17\x20\x43\x84\x47\x86\x45\x83\xcf\x84\xc9\x49\x00\x28\xc5\x26\xc4\xc6\xc3\xc5\xc2\xc1\xc7\xc5\xc5\x21\xc2\xc6\xfb\xbf\x1b\x94\x98\x7a\x80\x89\xa9\x85\x81\x81\xe7\x7f\x33\x88\xd7\xc2\xc0\xc8\x03\x4c\x23\x2c\xbb\x81\xe2\xca\x6c\x6c\x4f\x27\xf4\xbe\x58\x39\x8f\xe7\xe9\xec\x5d\xcf\xe6\x74\x42\x38\x6c\x4f\x76\x34\xbc\xe8\x58\xc3\xa8\x05\x82\x70\x15\xcf\x66\xae\x7b\xd9\x30\x0b\x9b\x0a\xce\x67\x2d\xf3\x9f\x76\x4f\x45\x17\x47\xd2\x8b\x62\xfa\xd3\x39\x1b\x80\x2c\x24\xa5\x38\x1c\xb1\x6b\x17\x35\x95\x91\xe6\x25\x22\x1d\xce\xf3\x62\xf1\x9c\xa7\x5d\x2b\x5f\xac\x98\xf1\xb4\x7f\x3b\x85\x9e\x27\xd1\x31\x44\xa8\xc0\x19\x32\xc8\x56\xf1\x3c\x6e\xda\xfa\x74\x49\xe7\xd3\xfe\xf5\x2f\x9a\xf7\x52\x66\x21\xa6\xdf\x48\x4c\x19\xc4\xb9\x98\xdc\xb8\x24\xcb\xd7\xc4\x26\x56\xf2\x63\x83\xed\xd9\xd4\x0d\x40\x8a\x40\xc8\xd0\xd1\xcb\x44\xe6\x28\x32\x62\x1b\xc3\xab\xcf\x57\xee\x7a\x3e\x73\x2f\x8c\x4f\xbc\xc3\x01\x01\x00\x00\xff\xff\xe3\xab\x8e\x9b\x26\x07\x00\x00"

func dicUniUnkDicBytes() ([]byte, error) {
	return bindataRead(
		_dicUniUnkDic,
		"dic/uni/unk.dic",
	)
}

func dicUniUnkDic() (*asset, error) {
	bytes, err := dicUniUnkDicBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "dic/uni/unk.dic", size: 1830, mode: os.FileMode(420), modTime: time.Unix(1452316537, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
