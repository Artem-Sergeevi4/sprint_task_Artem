package ftracker

import (
	"fmt"
	"math"
)

const (
	lenStep   = 0.65  // средняя длина шага в метрах.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return distance(action) / duration
}

func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	var distanceTr, speed, calories float64
	heightInM := height / cmInM // переменная для роста в метрах

	switch trainingType {
	case "Бег":
		distanceTr = distance(action)
		speed = meanSpeed(action, duration)
		calories = RunningSpentCalories(action, weight, duration)
	case "Ходьба":
		distanceTr = distance(action)
		speed = meanSpeed(action, duration)
		calories = WalkingSpentCalories(action, duration, weight, heightInM) // используем heightInM
	case "Плавание":
		distanceTr = distance(action)
		speed = swimmingMeanSpeed(lengthPool, countPool, duration)
		calories = SwimmingSpentCalories(lengthPool, countPool, duration, weight)

		fmt.Printf("Плавание: distance=%v, speed=%v, calories=%v\n", distanceTr, speed, calories)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration, distanceTr, speed, calories)
}

const (
	runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

func RunningSpentCalories(action int, weight, duration float64) float64 {
	return (runningCaloriesMeanSpeedMultiplier * meanSpeed(action, duration) * runningCaloriesMeanSpeedShift * weight / mInKm * duration * minInH)
}

const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

func WalkingSpentCalories(action int, duration, weight, heightInM float64) float64 {
	meanSpeedSec := meanSpeed(action, duration) * kmhInMsec
	return ((walkingCaloriesWeightMultiplier * weight) + (math.Pow(meanSpeedSec, 2) / heightInM * walkingSpeedHeightMultiplier * weight)) * duration * minInH
}

const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых калорий при плавании относительно скорости.
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании.
)

func swimmingDistance(lengthPool, countPool int) float64 {
	return float64(lengthPool*countPool) / mInKm // дистанция в километрах
}

func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return swimmingDistance(lengthPool, countPool) / duration
}

func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	return ((swimmingMeanSpeed(lengthPool, countPool, duration) + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration)
}
