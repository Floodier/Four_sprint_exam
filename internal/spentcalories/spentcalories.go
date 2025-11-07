package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	info := strings.Split(data, ",")
	if len(info) != 3 {
		return 0, "", 0, errors.New("неверный формат данных. Ожидается: шаги, вид активности, время")
	}
	stepsInfo := strings.TrimSpace(info[0])
	steps, err := strconv.Atoi(stepsInfo)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка в данных шагов %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов <= 0")
	}
	trainingType := strings.TrimSpace(info[1])
	if len(trainingType) == 0 {
		return 0, "", 0, fmt.Errorf("ошибка в виде активности")
	}
	durationInfo := strings.TrimSpace(info[2])
	duration, err := time.ParseDuration(durationInfo)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка в данных времени %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность времени <= 0")
	}
	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	trainingDistance := (stepLength * float64(steps)) / mInKm
	return trainingDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	hourTime := duration.Hours()
	if hourTime == 0 {
		return 0
	}
	speed := float64(distance(steps, height) / hourTime)

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	switch trainingType {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v\nДистанция: %.2f км\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			trainingType, duration, dist, speed, calories), nil
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v\nДистанция: %.2f км\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			trainingType, duration, dist, speed, calories), nil
	default:
		return fmt.Sprintf("неизвестный тип тренировки: %s", trainingType), nil
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("один из параметров <= 0")
	}
	runningTime := duration.Minutes()
	runCalories := (weight * meanSpeed(steps, height, duration) * runningTime) / minInH

	return runCalories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("один из параметров <= 0")
	}
	walkingTime := duration.Minutes()
	walkCalories := (weight * meanSpeed(steps, height, duration) * walkingTime) / minInH * walkingCaloriesCoefficient

	return walkCalories, nil
}
