package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type DBConn struct {
	Address  string `json:"address"`
	Login    string `json:"login"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

func GetDBConnectionData() (DBConn, error) {
	path, err := getPath("configs/db_conn.json")
	if err != nil {
		return DBConn{}, err
	}

	var c DBConn
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return DBConn{}, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return DBConn{}, err
	}

	return c, err
}

type ServerCfg struct {
	Port string `json:"port"`
}

func GetServerConfig() (ServerCfg, error) {
	path, err := getPath("configs/server_cfg.json")
	if err != nil {
		return ServerCfg{}, err
	}

	var c ServerCfg
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ServerCfg{}, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return ServerCfg{}, err
	}

	return c, err
}

func getPath(localPath string) (string, error) {
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	pathToPath := filepath.FromSlash(absPath + "/path.txt")

	ok, err := checkFileExistence(pathToPath)
	if err != nil {
		return "", err
	}

	if ok {
		path, err := readTextFile(pathToPath)
		if err != nil {
			return "", err
		}
		return strings.ReplaceAll(path, "\n", "") + "/" + localPath, nil
	}

	return filepath.FromSlash(absPath + "/" + localPath), nil
}

func checkFileExistence(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func readTextFile(path string) (string, error) {
	file, err := os.Open(path)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)

	var text string
	for scanner.Scan() {
		text += fmt.Sprintf("%v\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text, nil
}
