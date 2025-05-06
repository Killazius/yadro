package competitor

import "time"

// Competitor - структура для хранения данных об участнике гонки по биатлону
type Competitor struct {
	ID           int  // Уникальный идентификатор
	Registered   bool // Зарегистрирован ли участник
	Disqualified bool // Дисквалифицирован ли участник

	PlannedStart time.Time // Запланированное время начала заезда
	ActualStart  time.Time // Фактическое время старта участника
	EndTime      time.Time // Время завершения заезда участника

	StartSet bool // Установлено ли время старта
	Started  bool // Начал ли участник заезд
	Finished bool // Завершил ли участник заезд
	Hits     int  // Количество попаданий

	LapStartTimes        []time.Time     // Время начала каждого круга
	LapPenaltyStartTimes []time.Time     // Время начала каждого штрафного круга
	LapTimes             []time.Duration // Продолжительность каждого круга
	PenaltyTime          time.Duration   // Штрафное время участника
}
