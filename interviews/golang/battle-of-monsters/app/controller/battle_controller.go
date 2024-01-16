package controller

import (
	"battle-of-monsters/app/db"
	"battle-of-monsters/app/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListBattles(context *gin.Context) {
	var battle []models.Battle

	var result *gorm.DB

	if result = db.CONN.Find(&battle); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	log.Printf("Found %v battles", result.RowsAffected)
	context.JSON(http.StatusOK, &battle)
}

func StartBattle(context *gin.Context) {
	var battleRequest struct {
		MonsterA uint `json:"monsterA" binding:"required"`
		MonsterB uint `json:"monsterB" binding:"required"`
	}

	if err := context.BindJSON(&battleRequest); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	condition := db.CONN.Where("id IN ?", []uint{battleRequest.MonsterA, battleRequest.MonsterB})
	var monsters []models.Monster
	if result := condition.Find(&monsters); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	winner := simulateBattle(monsters)
	battle := models.Battle{
		MonsterAID: battleRequest.MonsterA,
		MonsterA:   monsters[0],
		MonsterBID: battleRequest.MonsterB,
		MonsterB:   monsters[1],
		Winner:     winner,
	}

	if result := db.CONN.Create(&battle); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	log.Printf("battle %v has been created", battle.ID)
	context.JSON(http.StatusCreated, &battle)
}

func DeleteBattle(context *gin.Context) {
	battleID := context.Param("battleID")
	var battle models.Battle

	if result := db.CONN.First(&battle, battleID); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result := db.CONN.Delete(&models.Battle{}, battleID); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	context.Status((http.StatusNoContent))
}

func getFirstAttackerIndex(monsters []models.Monster) int {
	if monsters[1].Speed > monsters[0].Speed {
		return 1
	}

	isSpeedEqual := monsters[1].Speed == monsters[0].Speed
	if !isSpeedEqual {
		return 0
	}

	if monsters[0].Attack > monsters[1].Attack {
		return 0
	}

	return 1
}

func simulateBattle(monsters []models.Monster) models.Monster {
	monstersCopy := []models.Monster{
		monsters[0],
		monsters[1],
	}

	nextAttacker := getFirstAttackerIndex(monsters)
	hasWinner := false
	for !hasWinner {
		hasWinner, nextAttacker = generateBattleAttack(&monstersCopy, nextAttacker)
		time.Sleep(time.Second)
	}

	return monsters[getOtherIndex(nextAttacker)]
}

func generateBattleAttack(monstersPointer *[]models.Monster, attackerIndex int) (bool, int) {
	monsters := *monstersPointer
	attacker := &monsters[attackerIndex]
	defender := &monsters[getOtherIndex(attackerIndex)]
	attackerDamage := (*attacker).Attack - (*defender).Defense
	if attackerDamage < 1 {
		attackerDamage = 1
	}

	(*defender).Hp -= attackerDamage
	if int((*defender).Hp) <= 0 {
		return true, getOtherIndex(attackerIndex)
	}

	return false, getOtherIndex(attackerIndex)
}

func getOtherIndex(index int) int {
	if index == 0 {
		return 1
	}

	return 0
}
