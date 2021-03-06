// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	core_kubernetes_io "github.com/solo-io/supergloo/pkg/api/external/kubernetes/core/v1"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

var (
	mConfigSnapshotIn  = stats.Int64("config.supergloo.solo.io/snap_emitter/snap_in", "The number of snapshots in", "1")
	mConfigSnapshotOut = stats.Int64("config.supergloo.solo.io/snap_emitter/snap_out", "The number of snapshots out", "1")

	configsnapshotInView = &view.View{
		Name:        "config.supergloo.solo.io_snap_emitter/snap_in",
		Measure:     mConfigSnapshotIn,
		Description: "The number of snapshots updates coming in",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	configsnapshotOutView = &view.View{
		Name:        "config.supergloo.solo.io/snap_emitter/snap_out",
		Measure:     mConfigSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(configsnapshotInView, configsnapshotOutView)
}

type ConfigEmitter interface {
	Register() error
	Mesh() MeshClient
	MeshIngress() MeshIngressClient
	MeshGroup() MeshGroupClient
	RoutingRule() RoutingRuleClient
	SecurityRule() SecurityRuleClient
	TlsSecret() TlsSecretClient
	Upstream() gloo_solo_io.UpstreamClient
	Pod() core_kubernetes_io.PodClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ConfigSnapshot, <-chan error, error)
}

func NewConfigEmitter(meshClient MeshClient, meshIngressClient MeshIngressClient, meshGroupClient MeshGroupClient, routingRuleClient RoutingRuleClient, securityRuleClient SecurityRuleClient, tlsSecretClient TlsSecretClient, upstreamClient gloo_solo_io.UpstreamClient, podClient core_kubernetes_io.PodClient) ConfigEmitter {
	return NewConfigEmitterWithEmit(meshClient, meshIngressClient, meshGroupClient, routingRuleClient, securityRuleClient, tlsSecretClient, upstreamClient, podClient, make(chan struct{}))
}

func NewConfigEmitterWithEmit(meshClient MeshClient, meshIngressClient MeshIngressClient, meshGroupClient MeshGroupClient, routingRuleClient RoutingRuleClient, securityRuleClient SecurityRuleClient, tlsSecretClient TlsSecretClient, upstreamClient gloo_solo_io.UpstreamClient, podClient core_kubernetes_io.PodClient, emit <-chan struct{}) ConfigEmitter {
	return &configEmitter{
		mesh:         meshClient,
		meshIngress:  meshIngressClient,
		meshGroup:    meshGroupClient,
		routingRule:  routingRuleClient,
		securityRule: securityRuleClient,
		tlsSecret:    tlsSecretClient,
		upstream:     upstreamClient,
		pod:          podClient,
		forceEmit:    emit,
	}
}

type configEmitter struct {
	forceEmit    <-chan struct{}
	mesh         MeshClient
	meshIngress  MeshIngressClient
	meshGroup    MeshGroupClient
	routingRule  RoutingRuleClient
	securityRule SecurityRuleClient
	tlsSecret    TlsSecretClient
	upstream     gloo_solo_io.UpstreamClient
	pod          core_kubernetes_io.PodClient
}

func (c *configEmitter) Register() error {
	if err := c.mesh.Register(); err != nil {
		return err
	}
	if err := c.meshIngress.Register(); err != nil {
		return err
	}
	if err := c.meshGroup.Register(); err != nil {
		return err
	}
	if err := c.routingRule.Register(); err != nil {
		return err
	}
	if err := c.securityRule.Register(); err != nil {
		return err
	}
	if err := c.tlsSecret.Register(); err != nil {
		return err
	}
	if err := c.upstream.Register(); err != nil {
		return err
	}
	if err := c.pod.Register(); err != nil {
		return err
	}
	return nil
}

func (c *configEmitter) Mesh() MeshClient {
	return c.mesh
}

func (c *configEmitter) MeshIngress() MeshIngressClient {
	return c.meshIngress
}

func (c *configEmitter) MeshGroup() MeshGroupClient {
	return c.meshGroup
}

func (c *configEmitter) RoutingRule() RoutingRuleClient {
	return c.routingRule
}

func (c *configEmitter) SecurityRule() SecurityRuleClient {
	return c.securityRule
}

func (c *configEmitter) TlsSecret() TlsSecretClient {
	return c.tlsSecret
}

func (c *configEmitter) Upstream() gloo_solo_io.UpstreamClient {
	return c.upstream
}

func (c *configEmitter) Pod() core_kubernetes_io.PodClient {
	return c.pod
}

func (c *configEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ConfigSnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	var done sync.WaitGroup
	ctx := opts.Ctx
	/* Create channel for Mesh */
	type meshListWithNamespace struct {
		list      MeshList
		namespace string
	}
	meshChan := make(chan meshListWithNamespace)
	/* Create channel for MeshIngress */
	type meshIngressListWithNamespace struct {
		list      MeshIngressList
		namespace string
	}
	meshIngressChan := make(chan meshIngressListWithNamespace)
	/* Create channel for MeshGroup */
	type meshGroupListWithNamespace struct {
		list      MeshGroupList
		namespace string
	}
	meshGroupChan := make(chan meshGroupListWithNamespace)
	/* Create channel for RoutingRule */
	type routingRuleListWithNamespace struct {
		list      RoutingRuleList
		namespace string
	}
	routingRuleChan := make(chan routingRuleListWithNamespace)
	/* Create channel for SecurityRule */
	type securityRuleListWithNamespace struct {
		list      SecurityRuleList
		namespace string
	}
	securityRuleChan := make(chan securityRuleListWithNamespace)
	/* Create channel for TlsSecret */
	type tlsSecretListWithNamespace struct {
		list      TlsSecretList
		namespace string
	}
	tlsSecretChan := make(chan tlsSecretListWithNamespace)
	/* Create channel for Upstream */
	type upstreamListWithNamespace struct {
		list      gloo_solo_io.UpstreamList
		namespace string
	}
	upstreamChan := make(chan upstreamListWithNamespace)
	/* Create channel for Pod */
	type podListWithNamespace struct {
		list      core_kubernetes_io.PodList
		namespace string
	}
	podChan := make(chan podListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup namespaced watch for Mesh */
		meshNamespacesChan, meshErrs, err := c.mesh.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Mesh watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, meshErrs, namespace+"-meshes")
		}(namespace)
		/* Setup namespaced watch for MeshIngress */
		meshIngressNamespacesChan, meshIngressErrs, err := c.meshIngress.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting MeshIngress watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, meshIngressErrs, namespace+"-meshingresses")
		}(namespace)
		/* Setup namespaced watch for MeshGroup */
		meshGroupNamespacesChan, meshGroupErrs, err := c.meshGroup.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting MeshGroup watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, meshGroupErrs, namespace+"-meshgroups")
		}(namespace)
		/* Setup namespaced watch for RoutingRule */
		routingRuleNamespacesChan, routingRuleErrs, err := c.routingRule.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting RoutingRule watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, routingRuleErrs, namespace+"-routingrules")
		}(namespace)
		/* Setup namespaced watch for SecurityRule */
		securityRuleNamespacesChan, securityRuleErrs, err := c.securityRule.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting SecurityRule watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, securityRuleErrs, namespace+"-securityrules")
		}(namespace)
		/* Setup namespaced watch for TlsSecret */
		tlsSecretNamespacesChan, tlsSecretErrs, err := c.tlsSecret.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting TlsSecret watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, tlsSecretErrs, namespace+"-tlssecrets")
		}(namespace)
		/* Setup namespaced watch for Upstream */
		upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Upstream watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-upstreams")
		}(namespace)
		/* Setup namespaced watch for Pod */
		podNamespacesChan, podErrs, err := c.pod.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Pod watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, podErrs, namespace+"-pods")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case meshList := <-meshNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case meshChan <- meshListWithNamespace{list: meshList, namespace: namespace}:
					}
				case meshIngressList := <-meshIngressNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case meshIngressChan <- meshIngressListWithNamespace{list: meshIngressList, namespace: namespace}:
					}
				case meshGroupList := <-meshGroupNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case meshGroupChan <- meshGroupListWithNamespace{list: meshGroupList, namespace: namespace}:
					}
				case routingRuleList := <-routingRuleNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case routingRuleChan <- routingRuleListWithNamespace{list: routingRuleList, namespace: namespace}:
					}
				case securityRuleList := <-securityRuleNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case securityRuleChan <- securityRuleListWithNamespace{list: securityRuleList, namespace: namespace}:
					}
				case tlsSecretList := <-tlsSecretNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case tlsSecretChan <- tlsSecretListWithNamespace{list: tlsSecretList, namespace: namespace}:
					}
				case upstreamList := <-upstreamNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
					}
				case podList := <-podNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case podChan <- podListWithNamespace{list: podList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *ConfigSnapshot)
	go func() {
		originalSnapshot := ConfigSnapshot{}
		currentSnapshot := originalSnapshot.Clone()
		timer := time.NewTicker(time.Second * 1)
		sync := func() {
			if originalSnapshot.Hash() == currentSnapshot.Hash() {
				return
			}

			stats.Record(ctx, mConfigSnapshotOut.M(1))
			originalSnapshot = currentSnapshot.Clone()
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}

		for {
			record := func() { stats.Record(ctx, mConfigSnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case meshNamespacedList := <-meshChan:
				record()

				namespace := meshNamespacedList.namespace
				meshList := meshNamespacedList.list

				currentSnapshot.Meshes[namespace] = meshList
			case meshIngressNamespacedList := <-meshIngressChan:
				record()

				namespace := meshIngressNamespacedList.namespace
				meshIngressList := meshIngressNamespacedList.list

				currentSnapshot.Meshingresses[namespace] = meshIngressList
			case meshGroupNamespacedList := <-meshGroupChan:
				record()

				namespace := meshGroupNamespacedList.namespace
				meshGroupList := meshGroupNamespacedList.list

				currentSnapshot.Meshgroups[namespace] = meshGroupList
			case routingRuleNamespacedList := <-routingRuleChan:
				record()

				namespace := routingRuleNamespacedList.namespace
				routingRuleList := routingRuleNamespacedList.list

				currentSnapshot.Routingrules[namespace] = routingRuleList
			case securityRuleNamespacedList := <-securityRuleChan:
				record()

				namespace := securityRuleNamespacedList.namespace
				securityRuleList := securityRuleNamespacedList.list

				currentSnapshot.Securityrules[namespace] = securityRuleList
			case tlsSecretNamespacedList := <-tlsSecretChan:
				record()

				namespace := tlsSecretNamespacedList.namespace
				tlsSecretList := tlsSecretNamespacedList.list

				currentSnapshot.Tlssecrets[namespace] = tlsSecretList
			case upstreamNamespacedList := <-upstreamChan:
				record()

				namespace := upstreamNamespacedList.namespace
				upstreamList := upstreamNamespacedList.list

				currentSnapshot.Upstreams[namespace] = upstreamList
			case podNamespacedList := <-podChan:
				record()

				namespace := podNamespacedList.namespace
				podList := podNamespacedList.list

				currentSnapshot.Pods[namespace] = podList
			}
		}
	}()
	return snapshots, errs, nil
}
