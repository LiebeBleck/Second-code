package main

import (
	"fmt"
	"time"
)

// Общие константы для вычислений
const (
	MInKm                            = 1000  // количество метров в одном километре
	MinInHours                       = 60    // количество минут в одном часе
	LenStep                          = 0.65  // длина одного шага
	CmInM                            = 100   // количество сантиметров в одном метре
	CaloriesMeanSpeedMultiplier      = 18    // множитель средней скорости бега
	CaloriesMeanSpeedShift           = 1.79  // коэффициент изменения средней скорости
	CaloriesWeightMultiplier         = 0.035 // коэффициент для веса при ходьбе
	CaloriesSpeedHeightMultiplier    = 0.029 // коэффициент для роста при ходьбе
	KmHInMsec                        = 0.278 // коэффициент для перевода км/ч в м/с
	SwimmingLenStep                  = 1.38  // длина одного гребка при плавании
	SwimmingCaloriesMeanSpeedShift   = 1.1   // коэффициент изменения средней скорости при плавании
	SwimmingCaloriesWeightMultiplier = 2     // множитель веса пользователя при плавании
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество повторов (шаги, гребки при плавании)
	LenStep      float64       // длина одного шага или гребка в м
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	durationInHours := t.Duration.Hours()
	if durationInHours == 0 {
		return 0 // Возвращаем 0, чтобы избежать деления на ноль
	}
	return t.distance() / durationInHours
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0 // Для конкретных тренировок метод будет переопределен
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo возвращает структуру InfoMessage с информацией о тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для всех видов тренировок
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возвращает количество потраченных килокалорий при беге.
func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) *
		r.Weight / MInKm * r.Duration.Hours() * MinInHours)
}

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64 // рост пользователя
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
func (w Walking) Calories() float64 {
	heightInMeters := w.Height / CmInM
	if heightInMeters == 0 {
		return 0 // Возвращаем 0, если рост равен нулю
	}
	speedInMSec := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight +
		(speedInMSec*speedInMSec/heightInMeters)*CaloriesSpeedHeightMultiplier*w.Weight) *
		w.Duration.Hours() * MinInHours)
}

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int // длина бассейна
	CountPool  int // количество пересечений бассейна
}

// meanSpeed возвращает среднюю скорость при плавании.
func (s Swimming) meanSpeed() float64 {
	durationInHours := s.Duration.Hours()
	if durationInHours == 0 {
		return 0 // Возвращаем 0, чтобы избежать деления на ноль
	}
	return float64(s.LengthPool*s.CountPool) / MInKm / durationInHours
}

// Calories возвращает количество калорий, потраченных при плавании.
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}

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
		CountPool:  5,
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
