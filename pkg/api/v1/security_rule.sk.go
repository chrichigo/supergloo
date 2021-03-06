// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewSecurityRule(namespace, name string) *SecurityRule {
	return &SecurityRule{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *SecurityRule) SetStatus(status core.Status) {
	r.Status = status
}

func (r *SecurityRule) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *SecurityRule) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.TargetMesh,
		r.SourceSelector,
		r.DestinationSelector,
		r.AllowedPaths,
		r.AllowedMethods,
	)
}

type SecurityRuleList []*SecurityRule
type SecurityrulesByNamespace map[string]SecurityRuleList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list SecurityRuleList) Find(namespace, name string) (*SecurityRule, error) {
	for _, securityRule := range list {
		if securityRule.Metadata.Name == name {
			if namespace == "" || securityRule.Metadata.Namespace == namespace {
				return securityRule, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find securityRule %v.%v", namespace, name)
}

func (list SecurityRuleList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, securityRule := range list {
		ress = append(ress, securityRule)
	}
	return ress
}

func (list SecurityRuleList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, securityRule := range list {
		ress = append(ress, securityRule)
	}
	return ress
}

func (list SecurityRuleList) Names() []string {
	var names []string
	for _, securityRule := range list {
		names = append(names, securityRule.Metadata.Name)
	}
	return names
}

func (list SecurityRuleList) NamespacesDotNames() []string {
	var names []string
	for _, securityRule := range list {
		names = append(names, securityRule.Metadata.Namespace+"."+securityRule.Metadata.Name)
	}
	return names
}

func (list SecurityRuleList) Sort() SecurityRuleList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list SecurityRuleList) Clone() SecurityRuleList {
	var securityRuleList SecurityRuleList
	for _, securityRule := range list {
		securityRuleList = append(securityRuleList, proto.Clone(securityRule).(*SecurityRule))
	}
	return securityRuleList
}

func (list SecurityRuleList) Each(f func(element *SecurityRule)) {
	for _, securityRule := range list {
		f(securityRule)
	}
}

func (list SecurityRuleList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *SecurityRule) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

func (byNamespace SecurityrulesByNamespace) Add(securityRule ...*SecurityRule) {
	for _, item := range securityRule {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace SecurityrulesByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace SecurityrulesByNamespace) List() SecurityRuleList {
	var list SecurityRuleList
	for _, securityRuleList := range byNamespace {
		list = append(list, securityRuleList...)
	}
	return list.Sort()
}

func (byNamespace SecurityrulesByNamespace) Clone() SecurityrulesByNamespace {
	cloned := make(SecurityrulesByNamespace)
	for ns, list := range byNamespace {
		cloned[ns] = list.Clone()
	}
	return cloned
}

var _ resources.Resource = &SecurityRule{}

// Kubernetes Adapter for SecurityRule

func (o *SecurityRule) GetObjectKind() schema.ObjectKind {
	t := SecurityRuleCrd.TypeMeta()
	return &t
}

func (o *SecurityRule) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*SecurityRule)
}

var SecurityRuleCrd = crd.NewCrd("supergloo.solo.io",
	"securityrules",
	"supergloo.solo.io",
	"v1",
	"SecurityRule",
	"sr",
	false,
	&SecurityRule{})
