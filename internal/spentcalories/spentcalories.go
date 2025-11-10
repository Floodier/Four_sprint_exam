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
		return 0, "", 0, errors.New("wrond data format. Expecting: steps, activity type, time")
	}
	stepsInfo := strings.TrimSpace(info[0])
	steps, err := strconv.Atoi(stepsInfo)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error in steps data: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps quantity less or equal zero")
	}
	trainingType := strings.TrimSpace(info[1])
	if len(trainingType) == 0 {
		return 0, "", 0, fmt.Errorf("error in activity type")
	}
	durationInfo := strings.TrimSpace(info[2])
	duration, err := time.ParseDuration(durationInfo)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error in time data: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("time is less or equal zero")
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
	speed := distance(steps, height) / hourTime

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	hours := duration.Hours()
	durationStr := fmt.Sprintf("%.2f ч.", hours)

	switch trainingType {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %s\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, durationStr, dist, speed, calories), nil
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %s\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, durationStr, dist, speed, calories), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("steps is %d - less or equal zero", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight is %.2f - less or equal zero", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("height is %.2f - less or equal zero", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration is %v - less or equal zero", duration)
	}

	runningTime := duration.Minutes()
	runCalories := (weight * meanSpeed(steps, height, duration) * runningTime) / minInH

	return runCalories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("steps is %d - less or equal zero", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight is %.2f - less or equal zero", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("height is %.2f - less or equal zero", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration is %v - less or equal zero", duration)
	}

	walkingTime := duration.Minutes()
	walkCalories := (weight * meanSpeed(steps, height, duration) * walkingTime) / minInH * walkingCaloriesCoefficient

	return walkCalories, nil
}
