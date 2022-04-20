package redis

/* 需要重构
var (
	crdb *ClusterClient
)

func init() {
	crdb = NewClusterClient(&ClusterOptions{
		ClusterOptions: redis.ClusterOptions{
			Addr: []string{"127.0.0.1:30021"},
		},
		Name: "test",
	})
	status, err := crdb.Ping(context.Background()).Result()
	fmt.Println(status, err)
}

func TestClusterClient_Set(t *testing.T) {
	k1 := generateKey()
	str, _ := crdb.SetEX(context.Background(), k1, "1000", time.Second*20).Result()
	if str != "OK" {
		t.Fatal("SetEX not valid", str)
	}
	str2, _ := crdb.Get(context.Background(), k1).Result()
	if str2 != "1000" {
		t.Fatal("Get not valid", str2)
	}
}

func TestClusterClient_Del(t *testing.T) {
	k1 := generateKey()
	crdb.SetEX(context.Background(), k1, "1000", time.Second*60).Result()
	delNum, err := crdb.Del(context.Background(), k1).Result()
	if delNum != 1 {
		t.Fatal("Del not valid", delNum, err)
	}
}

//TODO Exists 多个key会报错？
func TestClusterClient_Exists(t *testing.T) {
	k1 := generateKey()

	if num, err := crdb.Exists(context.Background(), k1).Result(); num != 0 {
		t.Fatal("got not valid", err)
	}
	crdb.SetEX(context.Background(), k1, "1", time.Second*60).Result()
	if num, err := crdb.Exists(context.Background(), k1).Result(); num != 1 {
		t.Fatal("got not valid", err)
	}
}

func TestClusterClient_ZRange(t *testing.T) {
	k1 := generateKey()
	crdb.ZAdd(context.Background(), k1, &redis.Z{Score: 1.55, Member: "key1"}).Result()
	crdb.ZAdd(context.Background(), k1, &redis.Z{Score: 1.56, Member: "key2"}).Result()
	crdb.Expire(context.Background(), k1, time.Minute)

	vals, err := crdb.ZRange(context.Background(), k1, 0, -1).Result()
	if len(vals) != 2 {
		t.Fatal("got not valid", err)
	}
	if delNum, _ := crdb.Del(context.Background(), k1).Result(); delNum < 1 {
		t.Fatal("del not valid", err)
	}
}
*/
