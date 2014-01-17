package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

//store configure parameters from conf file
var conf = make(map[string]string)

func InitConf() bool {
	f, err := os.Open("conf")
	defer f.Close()

	if err != nil {
		log.Panicf("failed to open conf file!")
		return false
	}

	freader := bufio.NewReader(f)
	line_num := 0
	for {
		line, err := freader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panicf("error when reading conf file!")
			return false
		}

		line_num++
		line = strings.TrimSpace(line)

		if strings.Contains(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			log.Panicf("line_num[%d]: %s format error!", line_num, line)
			return false
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])

		conf[k] = v
	}

	return true
}

//s is the key,d is default return value
func GetBool(k string, d bool) bool {
	v := GetString(k, "false")
	ret, err := strconv.ParseBool(v)
	if err != nil {
		return d
	}

	return ret
}

func GetInt32(k string, d int32) int32 {
	v := GetString(k, "")
	ret, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return d
	}
	return int32(ret)
}

func GetInt64(k string, d int64) int64 {
	v := GetString(k, "")
	ret, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return d
	}
	return ret
}

func GetFloat32(k string, d float32) float32 {
	v := GetString(k, "")
	ret, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return d
	}
	return float32(ret)
}

func GetString(k, d string) string {
	v, ok := conf[k]
	if !ok {
		return d
	}

	return v
}
