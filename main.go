package main

import (
	"fmt"
	"time"
)

// Общие константы
const (
	MInKm                            = 1000.0
	MinsInHour                       = 60.0
	LenStep                          = 0.65
	CmInM                            = 100.0
	CaloriesMeanSpeedMultiplier      = 18.0
	CaloriesMeanSpeedShift           = 1.79
	CaloriesWeightMultiplier         = 0.035
	CaloriesSpeedHeightMultiplier    = 0.029
	KmHInMsec                        = 0.278
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2.0
)

// Training общая структура для тренировок
type Training struct {
	TrainingType string        // Тип тренировки
	Action       int           // Количество шагов или гребков
	LenStep      float64       // Длина одного шага или гребка
	Duration     time.Duration // Длительность тренировки
	Weight       float64       // Вес пользователя
}

// distance возвращает дистанцию в километрах
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость в км/ч
func (t Training) meanSpeed() float64 {
	durationInHours := t.Duration.Hours()
	if durationInHours == 0 {
		return 0
	}
	return t.distance() / durationInHours
}

// Calories возвращает 0 (переопределяется в дочерних структурах)
func (t Training) Calories() float64 {
	return 0
}

// TrainingInfo формирует общую информацию о тренировке
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// InfoMessage структура для отображения информации о тренировке
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// String форматирует вывод информации о тренировке
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.1f минут\nДистанция: %.2f км\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f",
		i.TrainingType, i.Duration.Minutes(), i.Distance, i.Speed, i.Calories)
}

// CaloriesCalculator интерфейс для тренировок
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Running структура для бега
type Running struct {
	Training
}

// Calories рассчитывает калории для бега
func (r Running) Calories() float64 {
	return ((CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) *
		r.Weight / MInKm * r.Duration.Hours() * MinsInHour)
}

// TrainingInfo возвращает информацию о тренировке Бег
func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

// Walking структура для ходьбы
type Walking struct {
	Training
	Height float64 // Рост пользователя
}

// Calories рассчитывает калории для ходьбы
func (w Walking) Calories() float64 {
	heightInMeters := w.Height / CmInM
	if heightInMeters == 0 {
		return 0
	}
	speedInMSec := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight +
		(speedInMSec*speedInMSec/heightInMeters)*CaloriesSpeedHeightMultiplier*w.Weight) *
		w.Duration.Hours() * MinsInHour)
}

// TrainingInfo возвращает информацию о тренировке Ходьба
func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

// Swimming структура для плавания
type Swimming struct {
	Training
	LengthPool int // Длина бассейна
	CountPool  int // Количество пересечений
}

// meanSpeed возвращает среднюю скорость для плавания
func (s Swimming) meanSpeed() float64 {
	durationInHours := s.Duration.Hours()
	if durationInHours == 0 {
		return 0
	}
	return float64(s.LengthPool*s.CountPool) / MInKm / durationInHours
}

// Calories рассчитывает калории для плавания
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo переопределяет информацию о тренировке Плавание
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     float64(s.LengthPool*s.CountPool) / MInKm,
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData выводит информацию о тренировке, с учетом переопределенных калорий
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	// Переопределяем поле Calories, используя метод Calories() конкретной структуры
	info.Calories = training.Calories()
	return info.String()
}

// main демонстрация работы программы
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
