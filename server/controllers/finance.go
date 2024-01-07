package controllers

import (
	"net/http"
	"sort"
	"time"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type FinanceController struct {
	DB *gorm.DB
}

func NewFinanceController(db *gorm.DB, router *gin.RouterGroup) {
	financeController := FinanceController{DB: db}
	router.GET("/networth", financeController.GetNetWorthOverTime)
}

func (fc *FinanceController) GetNetWorthOverTime(context *gin.Context) {
	accounts, err := models.GetAllAccountsWithValues(fc.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	networth := rollup(5*time.Second, accounts)
	context.JSON(http.StatusOK, networth)
}

type NetWorthPoint struct {
	Date  time.Time       `json:"date"`
	Value decimal.Decimal `json:"value"`
}

// returns a sorted map by descending timestamps
func mapToSortedList(values map[time.Time]decimal.Decimal) []NetWorthPoint {
	times := []time.Time{}
	for k := range values {
		times = append(times, k)
	}

	sort.SliceStable(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	sorted := make([]NetWorthPoint, len(times))
	for i, ts := range times {
		sorted[i] = NetWorthPoint{
			Date:  ts,
			Value: values[ts],
		}
	}
	return sorted
}

func createTimeBuckets(accounts []models.Account, interval time.Duration) []time.Time {
	firstDate := accounts[0].Values[0].CreatedAt
	lastDate := accounts[0].Values[0].CreatedAt
	for _, account := range accounts {
		for _, value := range account.Values {
			if value.CreatedAt.Before(firstDate) {
				firstDate = value.CreatedAt
			}
			if value.CreatedAt.After(lastDate) {
				lastDate = value.CreatedAt
			}
		}
	}

	timeSpan := lastDate.Sub(firstDate)
	numBuckets := (timeSpan.Nanoseconds() / interval.Nanoseconds()) + 1
	buckets := []time.Time{}
	for i := int64(1); i <= numBuckets; i++ {
		next := time.Duration(i-1) * interval
		buckets = append(buckets, firstDate.Round(interval).Add(next))
	}

	return buckets
}

func fillBuckets(accounts []models.Account, values map[time.Time]decimal.Decimal, buckets []time.Time, interval time.Duration) {
	for _, account := range accounts {
		var bucketIndex int
		mostRecentValue := account.Values[0]
		for i, bucket := range buckets {
			if bucket == mostRecentValue.CreatedAt.Round(interval) {
				bucketIndex = i
			}
		}

		if bucketIndex < len(buckets)-1 {
			for i := len(buckets) - 1; i > bucketIndex; i-- {
				bucket := buckets[i]
				if account.Class == models.Asset {
					values[bucket] = values[bucket].Add(mostRecentValue.Value)
				} else {
					values[bucket] = values[bucket].Sub(mostRecentValue.Value)
				}
			}
		}
	}
}

func rollup(interval time.Duration, accounts []models.Account) []NetWorthPoint {
	buckets := createTimeBuckets(accounts, interval)
	values := make(map[time.Time]decimal.Decimal, len(buckets))
	fillBuckets(accounts, values, buckets, interval)

	for _, account := range accounts {
		for _, accountValue := range account.Values {
			ts := accountValue.CreatedAt.Round(interval)
			if account.Class == models.Asset {
				values[ts] = values[ts].Add(accountValue.Value)
			} else {
				values[ts] = values[ts].Sub(accountValue.Value)
			}
		}
	}
	return mapToSortedList(values)
}
