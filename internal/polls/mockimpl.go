package polls

import "fmt"

type MockRepo struct {
	polls []Model
}

func NewMockRepo() Repo {
	return &MockRepo{
		polls: []Model{
			{Id: 1, Text: "Что вы думаете о Go?", Type: FREE, AuthorID: 1, Cluster: 1},
			{Id: 2, Text: "Какой ваш любимый цвет?", Type: RADIO, AuthorID: 1, Cluster: 1, Answers: []string{"Красный", "Синий", "Зелёный"}},
			{Id: 3, Text: "Какая ваша любимая еда?", Type: RADIO, AuthorID: 2, Cluster: 2, Answers: []string{"Пицца", "Бургер", "Суши"}},
			{Id: 4, Text: "Какая ваша любимая музыка?", Type: RADIO, AuthorID: 2, Cluster: 2, Answers: []string{"Рок", "Поп", "Классика"}},
			{Id: 5, Text: "Где вы хотите провести отпуск?", Type: RADIO, AuthorID: 3, Cluster: 1, Answers: []string{"Моря", "Горы", "Города"}},
			{Id: 6, Text: "Какие ваши любимые жанры фильмов?", Type: RADIO, AuthorID: 3, Cluster: 3, Answers: []string{"Фантастика", "Драма", "Комедия"}},
			{Id: 7, Text: "Какой ваш любимый вид спорта?", Type: RADIO, AuthorID: 4, Cluster: 2, Answers: []string{"Футбол", "Теннис", "Баскетбол"}},
			{Id: 8, Text: "Какой язык программирования вы предпочитаете?", Type: RADIO, AuthorID: 4, Cluster: 1, Answers: []string{"Go", "Python", "JavaScript"}},
			{Id: 9, Text: "Что вы думаете о современных социальных сетях?", Type: FREE, AuthorID: 5, Cluster: 3},
			{Id: 10, Text: "Какой ваше отношение к искусственному интеллекту?", Type: FREE, AuthorID: 6, Cluster: 1},
			{Id: 11, Text: "Как часто вы используете интернет?", Type: RADIO, AuthorID: 7, Cluster: 2, Answers: []string{"Каждый день", "Каждую неделю", "Реже"}},
			{Id: 12, Text: "Какие книги вы любите читать?", Type: RADIO, AuthorID: 8, Cluster: 3, Answers: []string{"Романы", "Научную фантастику", "Нон-фикшн"}},
			{Id: 13, Text: "Какие фильмы вы предпочитаете?", Type: RADIO, AuthorID: 9, Cluster: 1, Answers: []string{"Драма", "Комедия", "Триллер"}},
			{Id: 14, Text: "Какие страны вам хотелось бы посетить?", Type: RADIO, AuthorID: 10, Cluster: 2, Answers: []string{"Япония", "Франция", "Италия"}},
			{Id: 15, Text: "Какое ваше мнение о здоровом питании?", Type: FREE, AuthorID: 11, Cluster: 3},
			{Id: 16, Text: "Какой вид спорта вам нравится больше всего?", Type: RADIO, AuthorID: 12, Cluster: 1, Answers: []string{"Футбол", "Хоккей", "Бокс"}},
			{Id: 17, Text: "Что вы думаете о долгосрочных отношениях?", Type: FREE, AuthorID: 13, Cluster: 2},
			{Id: 18, Text: "Какую профессию вы хотите выбрать?", Type: RADIO, AuthorID: 14, Cluster: 3, Answers: []string{"Программист", "Врач", "Учитель"}},
			{Id: 19, Text: "Какой праздник вы предпочитаете?", Type: RADIO, AuthorID: 15, Cluster: 1, Answers: []string{"Новый год", "День рождения", "Рождество"}},
			{Id: 20, Text: "Какой город в мире вы хотите посетить?", Type: RADIO, AuthorID: 16, Cluster: 2, Answers: []string{"Париж", "Нью-Йорк", "Лондон"}},
			{Id: 21, Text: "Какую музыку вы слушаете?", Type: RADIO, AuthorID: 17, Cluster: 3, Answers: []string{"Рок", "Поп", "Электроника"}},
			{Id: 22, Text: "Как вы относитесь к политике?", Type: FREE, AuthorID: 18, Cluster: 1},
			{Id: 23, Text: "Какие гаджеты вы используете?", Type: RADIO, AuthorID: 19, Cluster: 2, Answers: []string{"Смартфоны", "Ноутбуки", "Умные часы"}},
			{Id: 24, Text: "Что для вас важнее в жизни?", Type: FREE, AuthorID: 20, Cluster: 3},
			{Id: 25, Text: "Что вы думаете о будущем технологий?", Type: FREE, AuthorID: 21, Cluster: 1},
			{Id: 26, Text: "Как часто вы путешествуете?", Type: RADIO, AuthorID: 22, Cluster: 2, Answers: []string{"Раз в год", "Каждые несколько лет", "Очень редко"}},
			{Id: 27, Text: "Какой кофе вам нравится?", Type: RADIO, AuthorID: 23, Cluster: 3, Answers: []string{"Эспрессо", "Капучино", "Латте"}},
			{Id: 28, Text: "Какой фастфуд вы предпочитаете?", Type: RADIO, AuthorID: 24, Cluster: 1, Answers: []string{"Макдональдс", "KFC", "Бургер Кинг"}},
			{Id: 29, Text: "Какая ваша любимая погода?", Type: RADIO, AuthorID: 25, Cluster: 2, Answers: []string{"Жаркая", "Прохладная", "Холодная"}},
			{Id: 30, Text: "Какое ваше отношение к экологии?", Type: FREE, AuthorID: 26, Cluster: 3},
		},
	}
}

func (r *MockRepo) GetUncompletedPolls(pageIndex int, pageSize int, userId int64) ([]Model, error) {
	start := pageIndex * pageSize
	end := start + pageSize

	if start > len(r.polls) {
		return make([]Model, 0), nil
	}

	if end > len(r.polls) {
		end = len(r.polls)
	}

	return r.polls[start:end], nil
}

func (r *MockRepo) AddAnswer(userId int64, pollId int64, text string) error {
	fmt.Println("Add answer")
	return nil
}
