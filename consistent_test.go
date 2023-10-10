package consistent_hash

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	virtualNodeList := []int{100, 200, 300, 400}
	nodeNum := 10
	testCount := 1000000

	for _, v := range virtualNodeList {
		con := Consistent{}
		distributeMap := make(map[string]int64)

		for i := 1; i <= nodeNum; i++ {
			serverName := "172.168.0" + strconv.Itoa(i)
			con.AddNode(serverName, v)
			distributeMap[serverName] = 0
		}

		for i := 0; i < testCount; i++ {
			testName := "testName"
			serverName := con.GetNode(testName + strconv.Itoa(i))
			distributeMap[serverName] = distributeMap[serverName] + 1
		}
		var keys []string
		var values []float64
		for k, v := range distributeMap {
			keys = append(keys, k)
			values = append(values, float64(v))
		}
		sort.Strings(keys)
		fmt.Printf("####test %d nodes,every has %d virtual nodes ,%d test datas\n", nodeNum, v, testCount)
		for _, k := range keys {
			fmt.Printf("server address :%s data no:%d\n", k, distributeMap[k])
		}
		fmt.Printf("standard deviation\n:%f\n\n", getStandardDeviation(values))
	}
}

// get standard deviation
func getStandardDeviation(list []float64) float64 {
	var total float64
	for _, item := range list {
		total += item
	}
	// average no
	avg := total / float64(len(list))

	var dTotal float64
	for _, value := range list {
		dValue := value - avg
		dTotal += dValue * dValue
	}

	return math.Sqrt(dTotal / avg)
}
