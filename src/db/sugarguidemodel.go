package db

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

func SaveSugarGuidePlan(userId int, dietResult []byte, sportResult []byte, controlResult []byte) error {
	var dietPlan SugarGuideDietPlan
	var sportPlan SugarGuideSportPlan
	var controlPlan SugarGuideControlPlan
	if err := json.Unmarshal(dietResult, &dietPlan); err != nil {
		return err
	}
	if err := json.Unmarshal(sportResult, &sportPlan); err != nil {
		return err
	}
	if err := json.Unmarshal(controlResult, &controlPlan); err != nil {
		return err
	}
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return err
	}
	err = Transaction(func(db *gorm.DB) error {
		if err := db.Model(&user).Association("SugarGuideDietPlan").Clear().Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("SugarGuideDietPlan").Append(dietPlan).Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("SugarGuideSportPlan").Clear().Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("SugarGuideSportPlan").Append(sportPlan).Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("SugarGuideControlPlan").Clear().Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("SugarGuideControlPlan").Append(controlPlan).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func CheckWeeklyNewspaper(userId int) (bool, error) {
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return false, err
	}
	if mysqlDb.Model(&user).Association("SugarGuideDietPlan").Count() < 1 ||
		mysqlDb.Model(&user).Association("SugarGuideSportPlan").Count() < 1 ||
		mysqlDb.Model(&user).Association("SugarGuideControlPlan").Count() < 1 {
		return false, nil
	}
	return true, nil
}

func GetWeeklyNewspaper(userId int) (SugarGuideDietPlan, SugarGuideSportPlan, SugarGuideControlPlan, error) {
	var dietPlan SugarGuideDietPlan
	var sportPlan SugarGuideSportPlan
	var controlPlan SugarGuideControlPlan
	user, err := GetUserFromUserId(userId)
	if err != nil {
		return dietPlan, sportPlan, controlPlan, err
	}
	if err := mysqlDb.Model(&user).Association("SugarGuideDietPlan").Find(&dietPlan).Error; err != nil {
		return dietPlan, sportPlan, controlPlan, err
	}
	if err := mysqlDb.Model(&user).Association("SugarGuideSportPlan").Find(&sportPlan).Error; err != nil {
		return dietPlan, sportPlan, controlPlan, err
	}
	if err := mysqlDb.Model(&user).Association("SugarGuideControlPlan").Find(&controlPlan).Error; err != nil {
		return dietPlan, sportPlan, controlPlan, err
	}
	return dietPlan, sportPlan, controlPlan, nil
}
