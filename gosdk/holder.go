package sdkv2

import "sync"

type pluginsHolder struct {
	holder map[string]*Plugin
	mutex sync.RWMutex
}
var ph *pluginsHolder

func (d *pluginsHolder) get(key string) (*Plugin, bool) {
	d.mutex.RLock()

	defer d.mutex.RUnlock()
	val, exists := d.holder[key]
	return val, exists
}

func (d *pluginsHolder) add(key string, value *Plugin) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.holder[key] = value
}

func GetPluginHolder()*pluginsHolder{
	if ph==nil{
		ph = &pluginsHolder{}
		ph.holder = make(map[string]*Plugin)
	}

	return  ph
}