package model

import (
	"fmt"
)

// CreateImageTableIfNotExists Create ImageTable If Not Exists
func CreateImageTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS images(
		image_id INT NOT NULL AUTO_INCREMENT,
		image_url VARCHAR(256),
		content_id INT NOT NULL,
		PRIMARY KEY (image_id),
		FOREIGN KEY (content_id) REFERENCES contents(content_id)
		); `
	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create image table failed", err)
		return
	}
	fmt.Println("Create image table successed or it already exists")
}

// QueryMaxImageID 查询目前最大的 imageID
func QueryMaxImageID() int {
	var n int
	row := DB.QueryRow(" select max(image_id) from images")
	row.Scan(&n)
	return n
}
