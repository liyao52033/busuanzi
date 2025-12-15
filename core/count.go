package core

import (
	"context"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Incr(ctx, rk.SitePvKey).Result()
	pagePv, _ := _redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique).Result()

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, rk.SiteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, rk.PageUvKey).Result()

	// 时间维度统计 - 今日、本月
	_redis.Incr(ctx, rk.TodayPvKey)
	_redis.Incr(ctx, rk.MonthPvKey)
	_redis.PFAdd(ctx, rk.TodayUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.MonthUvKey, userIdentity)

	// 获取时间维度统计数据
	todayPv, _ := _redis.Get(ctx, rk.TodayPvKey).Int64()
	todayUv, _ := _redis.PFCount(ctx, rk.TodayUvKey).Result()
	yesterdayPv, _ := _redis.Get(ctx, rk.YesterdayPvKey).Int64()
	yesterdayUv, _ := _redis.PFCount(ctx, rk.YesterdayUvKey).Result()
	monthPv, _ := _redis.Get(ctx, rk.MonthPvKey).Int64()
	monthUv, _ := _redis.PFCount(ctx, rk.MonthUvKey).Result()

	// setExpire
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	// 为时间维度数据设置独立的过期时间
	dailyExpire := viper.GetDuration("bsz.dailyexpire") * time.Second
	monthlyExpire := viper.GetDuration("bsz.monthlyexpire") * time.Second
	go setExpireWithDuration(dailyExpire, rk.TodayPvKey, rk.TodayUvKey, rk.YesterdayPvKey, rk.YesterdayUvKey)
	go setExpireWithDuration(monthlyExpire, rk.MonthPvKey, rk.MonthUvKey)

	return Counts{
		SitePv:      sitePv,
		SiteUv:      siteUv,
		PagePv:      int64(pagePv),
		PageUv:      pageUv,
		TodayPv:     todayPv,
		TodayUv:     todayUv,
		YesterdayPv: yesterdayPv,
		YesterdayUv: yesterdayUv,
		MonthPv:     monthPv,
		MonthUv:     monthUv,
	}
}

// Put
// @description put data only
func Put(ctx context.Context, host string, path string, userIdentity string) {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	_redis.Incr(ctx, rk.SitePvKey)
	_redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique)

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)

	// 时间维度统计
	_redis.Incr(ctx, rk.TodayPvKey)
	_redis.Incr(ctx, rk.MonthPvKey)
	_redis.PFAdd(ctx, rk.TodayUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.MonthUvKey, userIdentity)

	// setExpire
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	// 为时间维度数据设置独立的过期时间
	dailyExpire := viper.GetDuration("bsz.dailyexpire") * time.Second
	monthlyExpire := viper.GetDuration("bsz.monthlyexpire") * time.Second
	go setExpireWithDuration(dailyExpire, rk.TodayPvKey, rk.TodayUvKey)
	go setExpireWithDuration(monthlyExpire, rk.MonthPvKey, rk.MonthUvKey)
}

// Get bsz counts
func Get(ctx context.Context, host string, path string) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Get(ctx, rk.SitePvKey).Int64()
	pagePv, _ := _redis.ZScore(ctx, rk.PagePvKey, rk.PathUnique).Result()

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, rk.SiteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, rk.PageUvKey).Result()

	// 获取时间维度统计数据
	todayPv, _ := _redis.Get(ctx, rk.TodayPvKey).Int64()
	todayUv, _ := _redis.PFCount(ctx, rk.TodayUvKey).Result()
	yesterdayPv, _ := _redis.Get(ctx, rk.YesterdayPvKey).Int64()
	yesterdayUv, _ := _redis.PFCount(ctx, rk.YesterdayUvKey).Result()
	monthPv, _ := _redis.Get(ctx, rk.MonthPvKey).Int64()
	monthUv, _ := _redis.PFCount(ctx, rk.MonthUvKey).Result()

	// setExpire
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	// 为时间维度数据设置独立的过期时间
	dailyExpire := viper.GetDuration("bsz.dailyexpire") * time.Second
	monthlyExpire := viper.GetDuration("bsz.monthlyexpire") * time.Second
	go setExpireWithDuration(dailyExpire, rk.TodayPvKey, rk.TodayUvKey, rk.YesterdayPvKey, rk.YesterdayUvKey)
	go setExpireWithDuration(monthlyExpire, rk.MonthPvKey, rk.MonthUvKey)

	return Counts{
		SitePv:      sitePv,
		SiteUv:      siteUv,
		PagePv:      int64(pagePv),
		PageUv:      pageUv,
		TodayPv:     todayPv,
		TodayUv:     todayUv,
		YesterdayPv: yesterdayPv,
		YesterdayUv: yesterdayUv,
		MonthPv:     monthPv,
		MonthUv:     monthUv,
	}
}

func getKeys(host string, path string) RKeys {
	var siteUnique = host
	var pathUnique = path

	// 兼容旧版本
	if viper.GetBool("bsz.pathStyle") == false {
		pathUnique = host + "&" + path
	}

	// encrypt
	switch viper.GetString("bsz.Encrypt") {
	case "MD516":
		siteUnique = tool.Md5(siteUnique)[8:24]
		pathUnique = tool.Md5(pathUnique)[8:24]
	case "MD532":
		siteUnique = tool.Md5(siteUnique)
		pathUnique = tool.Md5(pathUnique)
	default:
		siteUnique = tool.Md5(siteUnique)
		pathUnique = tool.Md5(pathUnique)
	}

	redisPrefix := viper.GetString("redis.prefix")

	siteUvKey := strings.Join([]string{redisPrefix, "site_uv", siteUnique}, ":")
	pageUvKey := strings.Join([]string{redisPrefix, "page_uv", siteUnique, pathUnique}, ":")

	sitePvKey := strings.Join([]string{redisPrefix, "site_pv", siteUnique}, ":")
	pagePvKey := strings.Join([]string{redisPrefix, "page_pv", siteUnique}, ":")

	// 时间维度Key
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	month := now.Format("2006-01")

	todayPvKey := strings.Join([]string{redisPrefix, "daily_pv", siteUnique, today}, ":")
	todayUvKey := strings.Join([]string{redisPrefix, "daily_uv", siteUnique, today}, ":")
	yesterdayPvKey := strings.Join([]string{redisPrefix, "daily_pv", siteUnique, yesterday}, ":")
	yesterdayUvKey := strings.Join([]string{redisPrefix, "daily_uv", siteUnique, yesterday}, ":")
	monthPvKey := strings.Join([]string{redisPrefix, "monthly_pv", siteUnique, month}, ":")
	monthUvKey := strings.Join([]string{redisPrefix, "monthly_uv", siteUnique, month}, ":")

	return RKeys{
		SitePvKey:      sitePvKey,
		SiteUvKey:      siteUvKey,
		PagePvKey:      pagePvKey,
		PageUvKey:      pageUvKey,
		SiteUnique:     siteUnique,
		PathUnique:     pathUnique,
		TodayPvKey:     todayPvKey,
		TodayUvKey:     todayUvKey,
		YesterdayPvKey: yesterdayPvKey,
		YesterdayUvKey: yesterdayUvKey,
		MonthPvKey:     monthPvKey,
		MonthUvKey:     monthUvKey,
	}
}
