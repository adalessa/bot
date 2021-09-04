package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Listeners struct {
	api *OpApi
}

func NewListeners(api *OpApi) *Listeners {
	return &Listeners{api: api}
}

func (l *Listeners) MessageListener() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, "!op") {
			return
		}
		command := strings.TrimPrefix(m.Content, "!op ")

		// s.ChannelMessageSend(m.ChannelID, "Hola mundo")
		switch {
		case strings.HasPrefix(command, "chapter"):
			l.getChapter(command, s, m)
		case strings.HasPrefix(command, "search"):
			l.search(command, s, m)
		}

	}
}

func (l *Listeners) search(command string, s *discordgo.Session, m *discordgo.MessageCreate) {
	search := strings.TrimPrefix(command, "search ")
	var entities []string

	elements := strings.Split(search, ", ")
	for _, element := range elements {
		entities = append(entities, strings.Trim(element, "\""))
	}

	encounter, err := l.api.SearchEncounter(entities)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error buscando el capitulo, sorry")
		fmt.Println(err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Se encontraron %d veces", encounter.Times))

	var output string

	for _, chapter := range encounter.Chapters {
		output += fmt.Sprintf("Chapter %d, Date: %s, title: %s\n", chapter.Number, chapter.ReleaseDate.Format("2006-01-02"), chapter.Title)
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{{Name: "encounters.txt", Reader: strings.NewReader(output)}},
	})
}
func (l *Listeners) getChapter(command string, s *discordgo.Session, m *discordgo.MessageCreate) {
	chapterNumberstr := strings.TrimPrefix(command, "chapter ")
	chapterNumber, err := strconv.Atoi(chapterNumberstr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "No entendi, usa por ejemplo `!op chapter 1`")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Buscando el chapter %d", chapterNumber))

	chapter, err := l.api.GetChapter(uint(chapterNumber))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error buscando, sorry")
		fmt.Println(err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Lo encontre")
	s.ChannelMessageSend(m.ChannelID, "Titulo: "+chapter.Title)
	formatReleaseDate := chapter.ReleaseDate.Format("2006-01-02")
	releaseDate, _ := time.Parse("2006-01-02", formatReleaseDate)
	diff := time.Now().Sub(releaseDate).Hours() / 24
	s.ChannelMessageSend(m.ChannelID, "Release Date: "+formatReleaseDate+" hace: "+fmt.Sprintf("%.0f", diff)+" Dias")
	s.ChannelMessageSend(m.ChannelID, chapter.ShortSummary)
	for _, link := range chapter.Links {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Link a %s: %s", link.Site, link.Url))
	}
}
