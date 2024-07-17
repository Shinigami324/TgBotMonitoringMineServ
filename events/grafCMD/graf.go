package grafcmd

import (
	"TgGraf/lib/er"
	"image/color"
	"path/filepath"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func NewgraphInHours() (err error) {
	// Создание нового графика
	p := plot.New()

	// Заголовок графика и осей
	p.Title.Text = "Мониторинг графика серверов"
	p.X.Label.Text = "Дни работаы сервера"
	p.Y.Label.Text = "Количество игроков онлайн"

	// Пример данных для графика
	points := plotter.XYs{
		{X: 1, Y: 1},
		{X: 2, Y: 4},
		{X: 3, Y: 10},
	}

	// Создание линий графика
	line, err := plotter.NewLine(points)
	if err != nil {
		return er.Wrap("Ошибка в создание графика", err)
	}
	line.Color = color.RGBA{B: 255, A: 255}

	// Добавление линий на график
	p.Add(line)

	filepath := filepath.Join("C:\\Users\\zduda\\Desktop\\g\\events\\grafCMD\\saveGraf", "graph.png")

	// Сохранение графика в файл
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath); err != nil {
		return er.Wrap("Ошибка в сохранении файла:", err)
	}

	return nil
}

func GetGraf() string {
	return "C:\\Users\\zduda\\Desktop\\g\\events\\grafCMD\\saveGraf\\graph.png"
}
