package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.ccom/hardius/go_final_project/pkg/functions"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)`

	res, err := db.Exec(query,
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))

	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	var tasks = []*Task{}

	var query string
	var rows *sql.Rows
	var err error
	if len(search) > 0 {

		t, err := time.Parse("02.01.2006", search)
		if err == nil {
			date := t.Format("20060102")
			query = `SELECT * FROM scheduler WHERE date = :date LIMIT :limit`
			rows, err = db.Query(query, sql.Named("date", date), sql.Named("limit", limit))
		} else {
			search := "%" + search + "%"
			query = `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
			rows, err = db.Query(query, sql.Named("search", search), sql.Named("limit", limit))
		}

	} else {
		query = `SELECT * FROM scheduler ORDER BY date LIMIT :limit`
		rows, err = db.Query(query, sql.Named("limit", limit))
	}

	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}

	return tasks, err
}

func GetTask(id string) (*Task, error) {
	task := Task{}

	query := `SELECT * FROM scheduler WHERE id = :id`
	err := db.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &Task{}, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	command := `UPDATE scheduler SET `
	arguments := `date = :date, title = :title, comment = :comment, repeat = :repeat `
	condition := `WHERE id = :id`

	query := command + arguments + condition
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`Incorrect id for updating task.`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"

	_, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func UpdateDate(task *Task) error {
	nextdate, err := functions.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	query := "UPDATE scheduler SET date = :date WHERE id = :id"
	res, err := db.Exec(query, sql.Named("date", nextdate), sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`Incorrect id for updating task.`)
	}
	return nil
}
