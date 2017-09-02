package install

import (
	sapi "github.com/appscode/stash/apis/stash"
	sapi_v1alpha1 "github.com/appscode/stash/apis/stash/v1alpha1"
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/pkg/api"
)

func init() {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  sapi.GroupName,
			VersionPreferenceOrder:     []string{sapi_v1alpha1.SchemeGroupVersion.Version},
			ImportPrefix:               "github.com/appscode/stash/apis/stash/v1alpha1",
			RootScopedKinds:            sets.NewString("CustomResourceDefinition"),
			AddInternalObjectsToScheme: sapi.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			sapi_v1alpha1.SchemeGroupVersion.Version: sapi_v1alpha1.AddToScheme,
		},
	).Announce(api.GroupFactoryRegistry).RegisterAndEnable(api.Registry, api.Scheme); err != nil {
		panic(err)
	}
}
