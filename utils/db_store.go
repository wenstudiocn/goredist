package utils

import (
	"bufio"
	"github.com/wenstudiocn/goredist/e"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

/**
Helper functions of data store. <Redis, MySQL ...>
*/

/**
Backup MySQL table rows into json file

@tableName: tableName, is not used for quering rows but as a part of json filename
@v: slice pointer used for gorm resultset, this is where to implement generic type function
@savePath: directory of where exported file located.
@version: version of exporting. used as a part of the json filename
*/
func BackupTable(tableName string, v interface{}, db *gorm.DB, savePath string, version string) error {
	value := reflect.ValueOf(v)
	if value.Type().Kind() != reflect.Ptr {
		return e.ErrParameters
	}
	rvalue := value.Elem()
	if rvalue.Kind() != reflect.Slice {
		return e.ErrParameters
	}
	// file
	filename := path.Join(savePath, fmt.Sprintf("%s.%s.json", tableName, version))
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	fmt.Printf("backup %s ...", tableName)

	w := bufio.NewWriter(fh)
	finished := 0
	for {
		err := db.Offset(finished).Limit(100).Find(v).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return err
		}
		if rvalue.Len() <= 0 {
			break
		}
		// write
		for i := 0; i < rvalue.Len(); i++ {
			user := rvalue.Index(i).Interface()
			data, err := json.Marshal(user)
			if err != nil {
				return err
			}
			w.Write(data)
			w.Write([]byte("\n"))
		}
		finished += rvalue.Len()

		fmt.Printf("...")
	}
	fmt.Printf("done\n")
	return w.Flush()
}

// TODO batch insertion
func RestoreTable(filefull string, v interface{}, db *gorm.DB) error {
	filename := filepath.Base(filefull)
	parts := strings.Split(filename, ".")
	if len(parts) != 3 {
		return e.ErrFormat
	}
	//tableName := parts[0]

	fh, err := os.OpenFile(filefull, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	reader := bufio.NewReader(fh)
	fmt.Printf("restoring %s ", filefull)
	count := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		err = json.Unmarshal(line, v)
		if err != nil {
			return err
		}
		err = db.Create(v).Error
		if err != nil {
			return err
		}

		count += 1
		if count % 100 == 0 {
			fmt.Printf("...")
		}
	}
	fmt.Println("Done.")
	return nil
}

/** Backup Database(or single table in it) using 'mysqldump' which must be pre-installed.

NOTICE: this function will block until done. You're able to delete 'cmd.Wait()' if you dont
care the result and want a async execution, and mostly which is OKay.
*/
func DumpDatabase(host, username, password, dbName, tableName, savePath, version string, port int) error {
	// command
	sPort := strconv.FormatInt(int64(port), 10)
	args := []string{"-h" + host, "-P" + sPort, "-u" + username, "-p" + password, dbName}
	if len(tableName) > 0 {
		args = append(args, tableName)
	} else {
		tableName = "all"
	}
	cmd := exec.Command("mysqldump", args...)
	// file
	filefull := path.Join(savePath, fmt.Sprintf("%s.%s.%s.sql", dbName, tableName, version))
	fh, err := os.OpenFile(filefull, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	cmd.Stdout = fh

	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
