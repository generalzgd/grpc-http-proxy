/**
 * @version: 1.0.0
 * @author: zgd: general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: data.go
 * @time: 2019-08-13 12:45
 */

package define

import (
	"sync"

	comm_libs `github.com/generalzgd/comm-libs`
	`hutte.zhanqi.tv/go/grpc-frame/common`
	"hutte.zhanqi.tv/go/link"
)

type UserLinks map[uint32]*link.Session // platform => link

type LinkCenter struct {
	lock sync.RWMutex
	data map[uint32]UserLinks // uid -> platform -> session
}

func NewLinkCenter() *LinkCenter {
	return &LinkCenter{
		data: map[uint32]UserLinks{},
	}
}

func is2n(v uint32) bool {
	if (v & (v - 1)) == 0 {
		return true
	}
	return false
}

// @return elder link -> need to close
// platform只允许1端，即2^n
func (p *LinkCenter) Put(uid uint32, platform uint32, newlink *link.Session) *link.Session {
	if !is2n(platform) {
		panic("platform value is not 2^n")
		return nil
	}
	// 所有移动端只有一个链接，也就是ios，android, ipad, tv互踢
	platform = p.filterMobilePlatform(platform)

	p.lock.Lock()
	defer p.lock.Unlock()

	var eld *link.Session
	userLinks, ok := p.data[uid]
	if !ok {
		userLinks = UserLinks{
			platform: newlink,
		}
	} else {
		eld = userLinks[platform]
		userLinks[platform] = newlink
	}
	p.data[uid] = userLinks
	return eld
}

func (p *LinkCenter) filterMobilePlatform(platform uint32) uint32 {
	t := uint32(0)
	if platform&1>0{
		t |= 1
	}
	if platform>=2{
		t |= 2
	}
	return t
}

// platform允许多端删除
func (p *LinkCenter) Delete(uid, platform uint32, expect *link.Session) []*link.Session {
	p.lock.Lock()
	defer p.lock.Unlock()

	// 所有移动端只有一个链接，也就是ios，android, ipad, tv互踢
	platform = p.filterMobilePlatform(platform)

	var li []*link.Session
	if links, ok := p.data[uid]; ok {
		for k, v := range links {
			if k&platform > 0 {
				if expect != nil && comm_libs.GetTarPointer(v) != comm_libs.GetTarPointer(expect) {
					continue
				}
				li = append(li, v)
				delete(links, k)
			}
		}
		p.data[uid] = links
	}
	return li
}

// platform允许多端
func (p *LinkCenter) Get(uid, platform uint32) []*link.Session {
	p.lock.RLock()
	defer p.lock.RUnlock()

	// 所有移动端只有一个链接，也就是ios，android, ipad, tv互踢
	platform = p.filterMobilePlatform(platform)

	var li []*link.Session
	if links, ok := p.data[uid]; ok {

		for k, v := range links {
			if k&platform > 0 {
				li = append(li, v)
			}
		}
	}
	return li
}

// ***********************************************************************************
// http客户端信息存储中心
type ClientInfoCenter struct {
	lock sync.RWMutex
	data map[string]*common.ClientConnInfo // guid(token)=>ClientConnInfo
}

func NewClientInfoCenter() *ClientInfoCenter {
	return &ClientInfoCenter{
		data: map[string]*common.ClientConnInfo{},
	}
}

func (p *ClientInfoCenter) Put(guid string, info *common.ClientConnInfo) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.data[guid] = info
}

func (p *ClientInfoCenter) Get(guid string) (*common.ClientConnInfo, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, ok := p.data[guid]
	if ok && v.Expired() {
		delete(p.data, guid)
		return nil, false
	}
	return v, ok
}

func (p *ClientInfoCenter) Delete(guid string, except *common.ClientConnInfo) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if except != nil {
		if v, ok := p.data[guid]; ok {
			if comm_libs.GetTarPointer(v) != comm_libs.GetTarPointer(except) {
				return
			}
		}
	}

	delete(p.data, guid)
}

func (p *ClientInfoCenter) Has(guid string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, ok := p.data[guid]
	if ok && v.Expired() {
		delete(p.data, guid)
		return false
	}
	return ok
}