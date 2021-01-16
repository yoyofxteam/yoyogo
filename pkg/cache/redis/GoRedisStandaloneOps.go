package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ctx = context.Background()
)

type GoRedisStandaloneOps struct {
	client redis.Cmdable
}

func NewStandaloneOps(options *Options) *GoRedisStandaloneOps {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &GoRedisStandaloneOps{client: client}
}

func (ops *GoRedisStandaloneOps) Ping() (string, error) {
	return ops.client.Ping(ctx).Result()
}

func (ops *GoRedisStandaloneOps) Close() error {
	return ops.client.(*redis.Client).Close()
}

// value ops

func (ops *GoRedisStandaloneOps) GetRange(key string, start int64, end int64) (string, error) {
	return ops.client.GetRange(ctx, key, start, end).Result()
}

func (ops *GoRedisStandaloneOps) StrLen(key string) (int64, error) {
	return ops.client.StrLen(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) Append(key string, value string) (int64, error) {
	return ops.client.Append(ctx, key, value).Result()
}

func (ops *GoRedisStandaloneOps) DeleteKey(keys ...string) (int64, error) {
	return ops.client.Del(ctx, keys...).Result()
}

func (ops *GoRedisStandaloneOps) Exists(key string) (bool, error) {
	n, e := ops.client.Exists(ctx, key).Result()
	return n > 0, e
}

func (ops *GoRedisStandaloneOps) SetExpire(key string, expiration time.Duration) (bool, error) {
	return ops.client.Expire(ctx, key, expiration).Result()
}

func (ops *GoRedisStandaloneOps) TTL(key string) (time.Duration, error) {
	return ops.client.TTL(ctx, key).Result()
}

// MultiSet is like Set but accepts multiple values:
//   - MSet("key1", "value1", "key2", "value2")
//   - MSet([]string{"key1", "value1", "key2", "value2"})
//   - MSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
func (ops *GoRedisStandaloneOps) MultiSet(values ...interface{}) error {
	return ops.client.MSet(ctx, values).Err()
}

func (ops *GoRedisStandaloneOps) SetValue(key string, value interface{}, expiration time.Duration) error {
	return ops.client.Set(ctx, key, value, expiration).Err()
}

func (ops *GoRedisStandaloneOps) Set(key string, value string, expiration time.Duration) error {
	return ops.client.Set(ctx, key, value, expiration).Err()
}

func (ops *GoRedisStandaloneOps) SetNX(key string, value interface{}) (bool, error) {
	return ops.client.SetNX(ctx, key, value, 0).Result()
}

func (ops *GoRedisStandaloneOps) GetValue(key string) ([]byte, error) {
	return ops.client.Get(ctx, key).Bytes()
}

func (ops *GoRedisStandaloneOps) Get(key string) (string, error) {
	return ops.client.Get(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) MultiGet(key ...string) ([]interface{}, error) {
	return ops.client.MGet(ctx, key...).Result()
}

func (ops *GoRedisStandaloneOps) IncrBy(key string, step int64) (int64, error) {
	return ops.client.IncrBy(ctx, key, step).Result()
}

func (ops *GoRedisStandaloneOps) RandomKey() (string, error) {
	return ops.client.RandomKey(ctx).Result()
}

//---------------------------------------------------------------------------------------------------
// list ops
func (ops *GoRedisStandaloneOps) LPop(key string) (string, error) {
	return ops.client.LPop(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) LIndex(key string, index int64) (string, error) {
	return ops.client.LIndex(ctx, key, index).Result()
}

func (ops *GoRedisStandaloneOps) LPush(key string, values ...interface{}) (int64, error) {
	return ops.client.LPush(ctx, key, values).Result()
}

//--------------------------------------------------------------------------------------------------
//geo ops

func (ops *GoRedisStandaloneOps) GeoAddArr(key string, geoLocation ...GeoPosition) int64 {
	var geoList = make([]*redis.GeoLocation, 0)
	for _, x := range geoLocation {
		geoEle := redis.GeoLocation{
			Longitude: x.Longitude,
			Latitude:  x.Latitude,
			Name:      x.Member,
		}
		geoList = append(geoList, &geoEle)
	}
	return ops.client.GeoAdd(ctx, key, geoList...).Val()
}

func (ops *GoRedisStandaloneOps) GeoPos(key string, members ...string) (error, []GeoPosition) {
	resList := ops.client.GeoPos(ctx, key, members...)
	if len(resList.Val()) == 0 {
		return errors.New("not find any geo info"), make([]GeoPosition, 0)
	}
	resGeoList := make([]GeoPosition, 0)
	resListVal := resList.Val()
	for i, x := range members {
		resValEle := resListVal[i]
		if resValEle != nil {
			resGeoList = append(resGeoList, GeoPosition{Longitude: resValEle.Longitude, Latitude: resValEle.Latitude, Member: x})
		}
	}
	return nil, resGeoList
}

func (ops *GoRedisStandaloneOps) GeoDist(key string, member1, member2 string, unit GeoUnit) (error, GeoDistInfo) {
	unitStr := getUnit(unit)
	if unitStr == "" {
		return errors.New("error unit"), GeoDistInfo{}
	}
	res := ops.client.GeoDist(ctx, key, member1, member2, unitStr).Val()
	return nil, GeoDistInfo{Unit: unit, Dist: res}
}

func (ops *GoRedisStandaloneOps) GeoRadius(key string, query GeoRadiusQuery) (error, []GeoPosition) {
	unitStr := getUnit(query.Unit)
	if unitStr == "" {
		return errors.New("error unit"), make([]GeoPosition, 0)
	}
	res := ops.client.GeoRadius(ctx, key, query.Longitude, query.Latitude, &redis.GeoRadiusQuery{
		Radius:      query.Radius,
		Unit:        unitStr,
		WithCoord:   query.WithCoord,
		WithDist:    query.WithDist,
		WithGeoHash: query.WithGeoHash,
		Count:       query.Count,
		Sort:        GetSort(query.Sort),
		Store:       query.Store,
		StoreDist:   query.StoreDist,
	})
	geoList := make([]GeoPosition, 0)
	for _, x := range res.Val() {
		geoList = append(geoList, GeoPosition{
			Member:    x.Name,
			Longitude: x.Longitude,
			Latitude:  x.Latitude,
			Dist:      x.Dist,
			GeoHash:   x.GeoHash,
			Unit:      query.Unit,
		})
	}
	return nil, geoList
}
func (ops *GoRedisStandaloneOps) GeoRadiusByMember(key string, member string, query GeoRadiusByMemberQuery) (error, []GeoPosition) {

	unitStr := getUnit(query.Unit)
	if unitStr == "" {
		return errors.New("error unit"), make([]GeoPosition, 0)
	}
	res := ops.client.GeoRadiusByMember(ctx, key, member, &redis.GeoRadiusQuery{
		Radius:      query.Radius,
		Unit:        unitStr,
		WithCoord:   query.WithCoord,
		WithDist:    query.WithDist,
		WithGeoHash: query.WithGeoHash,
		Count:       query.Count,
		Sort:        GetSort(query.Sort),
		Store:       query.Store,
		StoreDist:   query.StoreDist,
	})
	geoList := make([]GeoPosition, 0)
	for _, x := range res.Val() {
		geoList = append(geoList, GeoPosition{
			Member:    x.Name,
			Longitude: x.Longitude,
			Latitude:  x.Latitude,
			Dist:      x.Dist,
			GeoHash:   x.GeoHash,
		})
	}
	return nil, geoList
}

func (ops *GoRedisStandaloneOps) LRange(key string, start int64, end int64) ([]string, error) {
	return ops.client.LRange(ctx, key, start, end).Result()
}

func (ops *GoRedisStandaloneOps) LTrim(key string, start int64, end int64) error {
	return ops.client.LTrim(ctx, key, start, end).Err()
}

func (ops *GoRedisStandaloneOps) RPop(key string) (string, error) {
	return ops.client.RPop(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) RPush(key string, values ...interface{}) (int64, error) {
	return ops.client.RPush(ctx, key, values...).Result()
}

func (ops *GoRedisStandaloneOps) LSet(key string, index int64, value interface{}) error {
	return ops.client.LSet(ctx, key, index, value).Err()
}

func (ops *GoRedisStandaloneOps) LSize(key string) (int64, error) {
	return ops.client.LLen(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) LRemove(key string, count int64, value interface{}) (int64, error) {
	return ops.client.LRem(ctx, key, count, value).Result()
}

func (ops *GoRedisStandaloneOps) SAdd(key string, members ...interface{}) (int64, error) {
	return ops.client.SAdd(ctx, key, members...).Result()
}

func (ops *GoRedisStandaloneOps) SDiff(keys ...string) ([]string, error) {
	return ops.client.SDiff(ctx, keys...).Result()
}

func (ops *GoRedisStandaloneOps) SCard(key string) (int64, error) {
	return ops.client.SCard(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) SInter(keys ...string) ([]string, error) {
	return ops.client.SInter(ctx, keys...).Result()
}

func (ops *GoRedisStandaloneOps) SInterStore(destination string, keys ...string) (int64, error) {
	return ops.client.SInterStore(ctx, destination, keys...).Result()
}

func (ops *GoRedisStandaloneOps) SIsMember(key string, member interface{}) (bool, error) {
	return ops.client.SIsMember(ctx, key, member).Result()
}

func (ops *GoRedisStandaloneOps) SMembers(key string) ([]string, error) {
	return ops.client.SMembers(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) SMove(source string, destination string, member interface{}) (bool, error) {
	return ops.client.SMove(ctx, source, destination, member).Result()
}

func (ops *GoRedisStandaloneOps) SPop(key string) (string, error) {
	return ops.client.SPop(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) SRandMembers(key string, count int64) ([]string, error) {
	return ops.client.SRandMemberN(ctx, key, count).Result()
}

func (ops *GoRedisStandaloneOps) SRem(key string, members ...interface{}) (int64, error) {
	return ops.client.SRem(ctx, key, members...).Result()
}

func (ops *GoRedisStandaloneOps) SUnion(keys ...string) ([]string, error) {
	return ops.client.SUnion(ctx, keys...).Result()
}

func (ops *GoRedisStandaloneOps) SUnionStore(destination string, keys ...string) (int64, error) {
	return ops.client.SUnionStore(ctx, destination, keys...).Result()
}

func (ops *GoRedisStandaloneOps) HDel(key string, fields ...string) (int64, error) {
	return ops.client.HDel(ctx, key, fields...).Result()
}

func (ops *GoRedisStandaloneOps) HExists(key string, field string) (bool, error) {
	return ops.client.HExists(ctx, key, field).Result()
}

func (ops *GoRedisStandaloneOps) HGet(key string, field string) (string, error) {
	return ops.client.HGet(ctx, key, field).Result()
}

func (ops *GoRedisStandaloneOps) HGetInt64(key string, field string) (int64, error) {
	return ops.client.HGet(ctx, key, field).Int64()
}

func (ops *GoRedisStandaloneOps) HGetFloat64(key string, field string) (float64, error) {
	return ops.client.HGet(ctx, key, field).Float64()
}

func (ops *GoRedisStandaloneOps) HGetAll(key string) (map[string]string, error) {
	return ops.client.HGetAll(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) HIncrBy(key string, field string, increment int64) (int64, error) {
	return ops.client.HIncrBy(ctx, key, field, increment).Result()
}

func (ops *GoRedisStandaloneOps) HKeys(key string) ([]string, error) {
	return ops.client.HKeys(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) HLen(key string) (int64, error) {
	return ops.client.HLen(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) HMGet(key string, fields ...string) ([]interface{}, error) {
	return ops.client.HMGet(ctx, key, fields...).Result()
}

func (ops *GoRedisStandaloneOps) HSet(key string, field string, value interface{}) (int64, error) {
	return ops.client.HSet(ctx, key, field, value).Result()
}

func (ops *GoRedisStandaloneOps) HSetNX(key string, field string, value interface{}) (bool, error) {
	return ops.client.HSetNX(ctx, key, field, value).Result()
}

func (ops *GoRedisStandaloneOps) HVals(key string) ([]string, error) {
	return ops.client.HVals(ctx, key).Result()
}

func (ops *GoRedisStandaloneOps) ZAdd(key string, member ZMember) int64 {
	return ops.client.ZAdd(ctx, key, &redis.Z{Score: member.Score, Member: member.Member}).Val()
}

func (ops *GoRedisStandaloneOps) ZCard(key string) int64 {
	return ops.client.ZCard(ctx, key).Val()
}

func (ops *GoRedisStandaloneOps) ZCount(key, min, max string) int64 {
	return ops.client.ZCount(ctx, key, min, max).Val()
}

func (ops *GoRedisStandaloneOps) ZIncrby(key string, incr float64, member string) float64 {
	return ops.client.ZIncrBy(ctx, key, incr, member).Val()
}

func (ops *GoRedisStandaloneOps) ZInterStore(destination string, store []ZStore, arg ZStoreEnum) int64 {
	keyArr := make([]string, 0)
	weightArr := make([]float64, 0)
	for _, x := range store {
		keyArr = append(keyArr, x.Key)
		weightArr = append(weightArr, x.Weight)
	}
	var argStr = ""
	switch arg {
	case MAX:
		argStr = "MAX"
	case MIN:
		argStr = "MIN"
	case SUM:
		argStr = "SUM"
	default:
		argStr = "SUM"
	}
	return ops.client.ZInterStore(ctx, destination, &redis.ZStore{
		Keys:      keyArr,
		Weights:   weightArr,
		Aggregate: argStr,
	}).Val()
}

func (ops *GoRedisStandaloneOps) ZLexCount(key, min, max string) int64 {
	return ops.client.ZLexCount(ctx, key, min, max).Val()
}

func (ops *GoRedisStandaloneOps) ZRange(key string, start, stop int64) []string {
	return ops.client.ZRange(ctx, key, start, stop).Val()
}

func (ops *GoRedisStandaloneOps) ZRangeByLex(key, min, max string, offset int64, count int64) []string {
	return ops.client.ZRangeByLex(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Val()
}

func (ops *GoRedisStandaloneOps) ZRangeByScore(key, min, max string, offset int64, count int64) []string {
	return ops.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Val()
}

func (ops *GoRedisStandaloneOps) ZRank(key, member string) int64 {
	return ops.client.ZRank(ctx, key, member).Val()
}

func (ops *GoRedisStandaloneOps) ZRem(key string, member ...string) int64 {
	return ops.client.ZRem(ctx, key, member).Val()
}

func (ops *GoRedisStandaloneOps) ZRemRangeByLex(key, min, max string) int64 {
	return ops.client.ZRemRangeByLex(ctx, key, min, max).Val()
}

func (ops *GoRedisStandaloneOps) ZRemRangeByRank(key string, start, stop int64) int64 {
	return ops.client.ZRemRangeByRank(ctx, key, start, stop).Val()
}

func (ops *GoRedisStandaloneOps) ZRevRange(key string, start, stop int64) []string {
	return ops.client.ZRevRange(ctx, key, start, stop).Val()
}

func (ops *GoRedisStandaloneOps) ZRevRangeWithScores(key string, start, stop int64) ([]ZMember, error) {
	zArray, err := ops.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err == nil {
		zmArray := make([]ZMember, len(zArray))
		for i, z := range zArray {
			zmArray[i].Member = z.Member.(string)
			zmArray[i].Score = z.Score
		}
		return zmArray, err
	} else {
		return nil, err
	}
}

func (ops *GoRedisStandaloneOps) ZRevRank(key, member string) int64 {
	return ops.client.ZRevRank(ctx, key, member).Val()
}

func (ops *GoRedisStandaloneOps) ZScore(key, member string) float64 {
	return ops.client.ZScore(ctx, key, member).Val()
}
