package tests

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/last911/tools"
	"github.com/last911/tools/db"
	"testing"
	"time"
)

func TestMySQL(t *testing.T) {
	mysql, err := db.NewDB("mysql", "root:scnjl@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true")
	if err != nil {
		t.Fatal(err)
	}

	count, err := mysql.Count("SELECT COUNT(*) FROM test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("count:", count)

	// id, err := mysql.Insert("test", map[string]interface{}{
	// 	"name":     "scnjl",
	// 	"age":      38,
	// 	"dateline": tools.Now(),
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("ID:", id)

	count, err = mysql.Count("SELECT COUNT(*) FROM test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("count:", count)

	rows, err := mysql.FetchAll("SELECT comments, name, id, age FROM test")
	if err != nil {
		t.Fatal(err)
	}

	for _, row := range rows {
		for k, v := range row.(db.MapRow) {
			t.Logf("%s:%s", k, v)
		}
	}

	_, num, err := mysql.Execute("DELETE FROM test WHERE id<5")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("delete num:", num)

	test := Test{}
	row, err := mysql.FetchOne("SELECT comments, name, id, age, dateline FROM test WHERE id>10", &test)
	if err != nil {
		t.Fatal(err)
	}
	if tt, ok := row.(Test); ok {
		t.Log("Test:", tt)
	} else {
		t.Fatal("row:", row)
	}

	rows, err = mysql.FetchAll("SELECT comments, name, id, age, dateline FROM test WHERE id>10 and id<?", &Test{}, 90)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if tt, ok := row.(Test); ok {
			t.Log("Test:", tt)
		} else {
			t.Fatal("row:", row)
		}
	}
}

type Test struct {
	ID       int            `db:"id"`
	Name     string         `db:"name"`
	Age      int            `db:"age"`
	Dateline time.Time      `db:"dateline"`
	Comments sql.NullString `db:"comments"`
}
