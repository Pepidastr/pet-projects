package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

type User struct {
	id       int
	nickName string
	idUser   string
	points   int
}

func main() {
	db, err := sql.Open("sqlite", "gambling.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создаем сессию бота
	dg, err := discordgo.New("Bot " + "DISCORD_BOT_TOKEN")
	if err != nil {
		fmt.Println("Ошибка при создании сессии:", err)
		return
	}

	// Добавляем обработчик сообщений
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, db)
	})

	// Открываем соединение
	err = dg.Open()
	if err != nil {
		fmt.Println("Ошибка при открытии соединения:", err)
		return
	}
	defer dg.Close()

	fmt.Println("Бот запущен. Нажмите CTRL+C для завершения.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
	// Игнорируем сообщения от самого бота
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 10 && m.Content[0:9] == "!gambling" {
		user := User{
			id:       0,
			nickName: m.Author.GlobalName,
			idUser:   m.Author.ID,
			points:   0,
		}
		if m.Content[9:10] == " " && isDigit(m.Content[10:]) {
			num, err := strconv.Atoi(m.Content[10:])
			if err != nil {
				fmt.Println("Ошибка при преобразовании строки в число:", err)
				return
			}
			userGet, err := selectUser(db, user.idUser)
			if err != nil {
				if err == sql.ErrNoRows {
					insertUser(db, user)
					s.ChannelMessageSend(m.ChannelID, "Создан акк для -"+m.Author.Username)
					return
				} else {
					fmt.Println("Ошибка при выполнении запроса:", err)
					return
				}
			}
			if userGet.points < num {
				s.ChannelMessageSend(m.ChannelID, m.Author.Username+" недостаточно средств")
				return
			}
			if randomBool() {
				s.ChannelMessageSend(m.ChannelID, m.Author.Username+" выиграл "+strconv.Itoa(num)+" баллов")
				updatePointsOnUser(db, userGet.idUser, userGet.points+num)
			} else {
				s.ChannelMessageSend(m.ChannelID, m.Author.Username+" проиграл "+strconv.Itoa(num)+" баллов")
				updatePointsOnUser(db, userGet.idUser, userGet.points-num)
			}
			return
		} else {
			s.ChannelMessageSend(m.ChannelID, m.Author.Username+"! Некоректное использование команды")
		}
		return
	}

	if m.Content == "!info" {
		userGet, err := selectUser(db, m.Author.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				s.ChannelMessageSend(m.ChannelID, "Не существует аккаунта для - "+m.Author.Username)
				return
			} else {
				s.ChannelMessageSend(m.ChannelID, "Не известная ошибка, скажите дим пепу")
				return
			}
		}
		s.ChannelMessageSend(m.ChannelID, "Данные аккаунта -> "+userGet.nickName+"\n"+"ID аккаунта: "+userGet.idUser+"\n"+"Баллы: "+strconv.Itoa(userGet.points))
		return
	}

	if m.Content == "!mining" {
		userGet, err := selectUser(db, m.Author.ID)
		num := randomInt(3)
		if err != nil {
			if err == sql.ErrNoRows {
				s.ChannelMessageSend(m.ChannelID, "Не существует аккаунта для - "+m.Author.Username)
				return
			} else {
				s.ChannelMessageSend(m.ChannelID, "Не известная ошибка, скажите дим пепу")
				return
			}
		}
		s.ChannelMessageSend(m.ChannelID, userGet.nickName+" получил "+strconv.Itoa(num)+" баллов работая в шахте")
		updatePointsOnUser(db, userGet.idUser, userGet.points+num)
	}

	if m.Content == "!help" {
		s.ChannelMessageSend(m.ChannelID, `
		Это gambling bot, для создания аккаунта введите !gambling <число>
Доступныe команды:
	!gambling <число> - игра на удачу, выигрываете или проигрываете указанное число баллов шанс 50%
	<---------------------------------------------------------------------------------------------->
	!info - информация о аккаунте
	<---------------------------->
	!mining - безопасный способо добычи баллов, но более медленный
	<------------------------------------------------------------->`)
	}

}

func selectUser(db *sql.DB, userID string) (User, error) {
	var User User
	res := db.QueryRow("SELECT * FROM Bank WHERE idUser = :userID", sql.Named("userID", userID))
	err := res.Scan(&User.id, &User.nickName, &User.idUser, &User.points)
	if err != nil {
		return User, err
	}
	return User, nil
}

func insertUser(db *sql.DB, User User) error {
	_, err := db.Exec("INSERT INTO bank (nickName,idUser,points) VALUES (:nickName,:idUser,100)", sql.Named("nickName", User.nickName), sql.Named("idUser", User.idUser))
	if err != nil {
		return err
	}
	return nil
}

func updatePointsOnUser(db *sql.DB, userID string, points int) error {
	_, err := db.Exec("UPDATE Bank SET points = :points WHERE idUser = :idUser", sql.Named("points", points), sql.Named("idUser", userID))
	if err != nil {
		return err
	}
	return nil
}

func randomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}

func randomInt(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func isDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
