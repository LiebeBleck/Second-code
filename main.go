package main

import (
	"fmt"
	"time"
)

// Общие константы для вычислений
const (
	MInKm                            = 1000.0 // количество метров в одном километре
	MinsInHour                       = 60.0   // количество минут в одном часе
	LenStep                          = 0.65   // длина одного шага в метрах
	CmInM                            = 100.0  // количество сантиметров в одном метре
	CaloriesMeanSpeedMultiplier      = 18.0   // множитель средней скорости бега
	CaloriesMeanSpeedShift           = 1.79   // коэффициент изменения средней скорости
	CaloriesWeightMultiplier         = 0.035  // коэффициент для веса при ходьбе
	CaloriesSpeedHeightMultiplier    = 0.029  // коэффициент для роста при ходьбе
	KmHInMsec                        = 0.278  // коэффициент для перевода км/ч в м/с
	SwimmingLenStep                  = 1.38   // длина одного гребка при плавании
	SwimmingCaloriesMeanSpeedShift   = 1.1    // коэффициент изменения средней скорости при плавании
	SwimmingCaloriesWeightMultiplier = 2.0    // множитель веса пользователя при плавании
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // Тип тренировки
	Action       int           // Количество шагов или гребков
	LenStep      float64       // Длина одного шага или гребка
	Duration     time.Duration // Длительность тренировки
	Weight       float64       // Вес пользователя
}

// distance возвращает дистанцию в километрах, которую преодолел пользователь
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость в км/ч
func (t Training) meanSpeed() float64 {
	durationInHours := t.Duration.Hours()
	if durationInHours == 0 {
		return 0 // Избегаем деления на ноль
	}
	return t.distance() / durationInHours
}

// Calories возвращает количество потраченных калорий (переопределяется в дочерних структурах)
func (t Training) Calories() float64 {
	return 0
}

// TrainingInfo возвращает структуру InfoMessage с информацией о тренировке
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// InfoMessage содержит информацию о тренировке
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// String возвращает строку с информацией о тренировке
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.1f минут\nДистанция: %.2f км\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f",
		i.TrainingType, i.Duration.Minutes(), i.Distance, i.Speed, i.Calories)
}

// CaloriesCalculator интерфейс для всех видов тренировок
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Running структура для тренировки Бег
type Running struct {
	Training
}

// Calories рассчитывает калории для бега
func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) *
		r.Weight / MInKm * r.Duration.Hours() * MinsInHour)
}

// Walking структура для тренировки Ходьба
type Walking struct {
	Training
	Height float64 // Рост пользователя в сантиметрах
}

// Calories рассчитывает калории для ходьбы
func (w Walking) Calories() float64 {
	heightInMeters := w.Height / CmInM
	if heightInMeters == 0 {
		return 0 // Избегаем деления на ноль
	}
	speedInMSec := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight +
		(speedInMSec*speedInMSec/heightInMeters)*CaloriesSpeedHeightMultiplier*w.Weight) *
		w.Duration.Hours() * MinsInHour)
}

// Swimming структура для тренировки Плавание
type Swimming struct {
	Training
	LengthPool int // Длина бассейна в метрах
	CountPool  int // Количество пересечений бассейна
}

// meanSpeed возвращает среднюю скорость в км/ч
func (s Swimming) meanSpeed() float64 {
	durationInHours := s.Duration.Hours()
	if durationInHours == 0 {
		return 0 // Избегаем деления на ноль
	}
	return float64(s.LengthPool*s.CountPool) / MInKm / durationInHours
}

// Calories рассчитывает калории для плавания
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo переопределяет метод для получения информации о плавании
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     float64(s.LengthPool*s.CountPool) / MInKm,
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData получает информацию о тренировке
func ReadData(training CaloriesCalculator) string {
	return training.TrainingInfo().String()
}

// main пример работы программы
func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  40,
	}

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(swimming))
	fmt.Println(ReadData(walking))
	fmt.Println(ReadData(running))
}
