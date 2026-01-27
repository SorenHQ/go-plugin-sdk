package sdkv2

import "sync"

type jobToReqMap struct {
	holder map[string]string
	mutex sync.RWMutex
}
var sj *jobToReqMap // Space Jobs

func (d *jobToReqMap) Get(key string) (string, bool) {
	d.mutex.RLock()

	defer d.mutex.RUnlock()
	val, exists := d.holder[key]
	return val, exists
}

func (d *jobToReqMap) Add(jobId string, entityId string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.holder[jobId] = entityId
}

func (d *jobToReqMap) Delete(jobId string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.holder,jobId)
}

func GetjobsHolder()*jobToReqMap{
	if sj==nil{
		sj = &jobToReqMap{}
		sj.holder = make(map[string]string)
	}

	return  sj
}