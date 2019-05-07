package watcher

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/deislabs/smi-sdk-go/gen/apis/trafficsplit/v1beta1"
	"github.com/linkerd/linkerd2/controller/k8s"
	pkgK8s "github.com/linkerd/linkerd2/pkg/k8s"
	logging "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
)

type (
	TrafficSplitWatcher struct {
		endpoints  *EndpointsWatcher
		services   map[string]*apexServicePublisher
		servicesMu sync.RWMutex
		log        *logging.Entry
		k8sAPI     *k8s.API
	}

	apexServicePublisher struct {
		endpoints     *EndpointsWatcher
		authority     string
		pods          PodSet
		exists        bool
		listeners     []EndpointUpdateListener
		leafListeners map[string]EndpointUpdateListener
		log           *logging.Entry
		mutex         sync.Mutex
	}

	leafListener struct {
		pods   PodSet
		weight string
		parent *apexServicePublisher
		log    *logging.Entry
	}
)

func NewTrafficSplitWatcher(endpoints *EndpointsWatcher, k8sAPI *k8s.API, log *logging.Entry) *TrafficSplitWatcher {
	watcher := &TrafficSplitWatcher{
		endpoints: endpoints,
		services:  make(map[string]*apexServicePublisher),
		log:       log.WithField("component", "traffic-split-watcher"),
		k8sAPI:    k8sAPI,
	}

	k8sAPI.TS().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    watcher.addTrafficSplit,
			UpdateFunc: watcher.updateTrafficSplit,
			DeleteFunc: watcher.deleteTrafficSplit,
		},
	)

	return watcher
}

///////////////////////////
/// TrafficSplitWatcher ///
///////////////////////////

func (tsw *TrafficSplitWatcher) Subscribe(authority string, listener EndpointUpdateListener) error {
	tsw.servicesMu.Lock()
	asp, ok := tsw.services[authority]
	if !ok {
		service, _, err := GetServiceAndPort(authority)
		if err != nil {
			return err
		}
		trafficSplits, err := tsw.k8sAPI.TS().Lister().TrafficSplits(service.Namespace).List(labels.Everything())
		if err != nil {
			return err
		}
		var trafficSplit *v1beta1.TrafficSplit
		for _, ts := range trafficSplits {
			// Use the first TrafficSplit we find for this service.  If there are more than one
			// TrafficSplit resources for this service, disregard the others.
			if ts.Spec.Service == service.Name {
				trafficSplit = ts
				break
			}
		}
		asp = tsw.newApexServicePublisher(trafficSplit, authority)
		tsw.services[authority] = asp
	}
	tsw.servicesMu.Unlock()

	asp.subscribe(listener)
	return nil
}

func (tsw *TrafficSplitWatcher) Unsubscribe(authority string, listener EndpointUpdateListener) {
	tsw.servicesMu.RLock()
	asp, ok := tsw.services[authority]
	tsw.servicesMu.RUnlock()
	if ok {
		asp.unsubscribe(listener)
	}
}

func (tsw *TrafficSplitWatcher) addTrafficSplit(obj interface{}) {
	trafficSplit := obj.(*v1beta1.TrafficSplit)

	tsw.servicesMu.RLock()
	defer tsw.servicesMu.RUnlock()

	for authority, asp := range tsw.services {
		service, _, err := GetServiceAndPort(authority)
		if err != nil {
			tsw.log.Errorf("Malformed authority in services map: %s", authority)
			continue
		}
		if service.Name == trafficSplit.Spec.Service && service.Namespace == trafficSplit.Namespace {
			asp.setTrafficSplit(trafficSplit)
			break
		}
	}
}

func (tsw *TrafficSplitWatcher) deleteTrafficSplit(obj interface{}) {
	trafficSplit := obj.(*v1beta1.TrafficSplit)

	tsw.servicesMu.RLock()
	defer tsw.servicesMu.RUnlock()
	for authority, asp := range tsw.services {
		service, _, err := GetServiceAndPort(authority)
		if err != nil {
			tsw.log.Errorf("Malformed authority in services map: %s", authority)
			continue
		}
		if service.Name == trafficSplit.Spec.Service && service.Namespace == trafficSplit.Namespace {
			asp.setTrafficSplit(nil)
			break
		}
	}
}

func (tsw *TrafficSplitWatcher) updateTrafficSplit(old interface{}, new interface{}) {
	tsw.addTrafficSplit(new)
}

func (tsw *TrafficSplitWatcher) newApexServicePublisher(trafficSplit *v1beta1.TrafficSplit, authority string) *apexServicePublisher {

	asp := &apexServicePublisher{
		endpoints:     tsw.endpoints,
		pods:          make(PodSet),
		authority:     authority,
		listeners:     []EndpointUpdateListener{},
		leafListeners: make(map[string]EndpointUpdateListener),
		log:           tsw.log.WithField("authority", authority),
	}
	asp.setTrafficSplit(trafficSplit)
	return asp
}

////////////////////////////
/// apexServicePublisher ///
////////////////////////////

func (asp *apexServicePublisher) subscribe(listener EndpointUpdateListener) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()
	asp.listeners = append(asp.listeners, listener)
	if asp.exists {
		if len(asp.pods) > 0 {
			listener.Add(asp.pods)
			return
		}
		listener.NoEndpoints(true)
		return
	}
	listener.NoEndpoints(false)
}

func (asp *apexServicePublisher) unsubscribe(listener EndpointUpdateListener) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()

	for i, item := range asp.listeners {
		if listener == item {
			n := len(asp.listeners)
			asp.listeners[i] = asp.listeners[n-1]
			asp.listeners[n-1] = nil
			asp.listeners = asp.listeners[:n-1]
			break
		}
	}
}

func (asp *apexServicePublisher) setTrafficSplit(trafficSplit *v1beta1.TrafficSplit) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()

	for leaf, listener := range asp.leafListeners {
		asp.endpoints.Unsubscribe(leaf, listener)
	}

	// We temporarily release the lock so that we can accept an update that
	// comes synchronously when we subscribe.
	asp.mutex.Unlock()
	asp.Remove(asp.pods)
	asp.mutex.Lock()

	if trafficSplit == nil {
		// We temporarily release the lock so that we can accept an update that
		// comes synchronously when we subscribe.
		asp.mutex.Unlock()
		asp.endpoints.Subscribe(asp.authority, asp)
		asp.mutex.Lock()
		asp.leafListeners = map[string]EndpointUpdateListener{asp.authority: asp}
		return
	}

	asp.leafListeners = make(map[string]EndpointUpdateListener)
	_, port, err := getHostAndPort(asp.authority)
	if err != nil {
		asp.log.Errorf("Invalid authority in apex service publisher [%s]: %s", asp.authority, err)
		return
	}
	for _, backend := range trafficSplit.Spec.Backends {
		leafAuthority := fmt.Sprintf("%s.%s.svc.cluster.local:%d", backend.Service, trafficSplit.Namespace, port)
		listener := asp.newLeafListener(backend.Weight, leafAuthority)
		// We temporarily release the lock so that we can accept an update that
		// comes synchronously when we subscribe.
		asp.mutex.Unlock()
		asp.endpoints.Subscribe(leafAuthority, listener)
		asp.mutex.Lock()
		asp.leafListeners[leafAuthority] = listener
	}
}

// Fanout update to listeners.
func (asp *apexServicePublisher) Add(set PodSet) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()

	asp.exists = true
	for id, address := range set {
		asp.pods[id] = address
	}
	for _, listener := range asp.listeners {
		listener.Add(set)
	}
}

// Fanout update to listeners.
func (asp *apexServicePublisher) Remove(set PodSet) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()

	for _, listener := range asp.listeners {
		listener.Remove(set)
	}

	for id := range set {
		delete(asp.pods, id)
	}
}

// Fanout update to listeners.
func (asp *apexServicePublisher) NoEndpoints(exists bool) {
	asp.mutex.Lock()
	defer asp.mutex.Unlock()

	asp.pods = make(PodSet)
	asp.exists = exists
	for _, listener := range asp.listeners {
		listener.NoEndpoints(exists)
	}
}

func (asp *apexServicePublisher) newLeafListener(weight string, leafAuthority string) *leafListener {
	return &leafListener{
		parent: asp,
		pods:   make(PodSet),
		weight: weight,
		log:    asp.log.WithField("leaf-service", leafAuthority),
	}
}

////////////////////
/// leafListener ///
////////////////////

func (ll *leafListener) Add(set PodSet) {
	for id, address := range set {
		address.Pod.Annotations[pkgK8s.ProxyPodWeightAnnotation] = ll.weight
		ll.pods[id] = address
	}
	ll.parent.Add(set)
}

func (ll *leafListener) Remove(set PodSet) {
	for id := range set {
		delete(ll.pods, id)
	}
	ll.parent.Remove(set)
}

func (ll *leafListener) NoEndpoints(exists bool) {
	ll.parent.Remove(ll.pods)
	ll.pods = make(PodSet)
}
