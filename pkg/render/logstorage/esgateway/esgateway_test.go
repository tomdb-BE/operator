// Copyright (c) 2021 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package esgateway

import (
	"fmt"

	"github.com/tigera/operator/pkg/render/kubecontrollers"

	relasticsearch "github.com/tigera/operator/pkg/render/common/elasticsearch"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	operatorv1 "github.com/tigera/operator/api/v1"
	"github.com/tigera/operator/pkg/render"
	rmeta "github.com/tigera/operator/pkg/render/common/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type resourceTestObj struct {
	name string
	ns   string
	typ  runtime.Object
	f    func(resource runtime.Object)
}

var _ = Describe("ES Gateway rendering tests", func() {
	Context("ES Gateway deployment", func() {
		var installation *operatorv1.InstallationSpec
		clusterDomain := "cluster.local"

		BeforeEach(func() {
			installation = &operatorv1.InstallationSpec{
				KubernetesProvider: operatorv1.ProviderNone,
				Registry:           "testregistry.com/",
			}
		})

		It("should render an ES Gateway deployment and all supporting resources", func() {
			expectedResources := []resourceTestObj{
				{relasticsearch.PublicCertSecret, rmeta.OperatorNamespace(), &corev1.Secret{}, nil},
				{render.TigeraElasticsearchCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{relasticsearch.PublicCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{render.KibanaInternalCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersUserSecret, rmeta.OperatorNamespace(), &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersVerificationUserSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersSecureUserSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{ServiceName, render.ElasticsearchNamespace, &corev1.Service{}, nil},
				{RoleName, render.ElasticsearchNamespace, &rbacv1.Role{}, nil},
				{RoleName, render.ElasticsearchNamespace, &rbacv1.RoleBinding{}, nil},
				{ServiceAccountName, render.ElasticsearchNamespace, &corev1.ServiceAccount{}, nil},
				{DeploymentName, render.ElasticsearchNamespace, &appsv1.Deployment{}, nil},
			}

			component := EsGateway(&Config{
				installation,
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: "tigera-pull-secret"}},
				},
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: render.TigeraElasticsearchCertSecret, Namespace: rmeta.OperatorNamespace()}},
					{ObjectMeta: metav1.ObjectMeta{Name: relasticsearch.PublicCertSecret, Namespace: rmeta.OperatorNamespace()}},
				},
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersUserSecret, Namespace: rmeta.OperatorNamespace()}},
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersVerificationUserSecret, Namespace: render.ElasticsearchNamespace}},
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersSecureUserSecret, Namespace: render.ElasticsearchNamespace}},
				},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: render.KibanaInternalCertSecret, Namespace: rmeta.OperatorNamespace()}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: relasticsearch.InternalCertSecret, Namespace: render.ElasticsearchNamespace}},
				clusterDomain, "elastic",
			})

			createResources, _ := component.Objects()
			compareResources(createResources, expectedResources)
		})

		It("should render an ES Gateway deployment and all supporting resources when CertificateManagement is enabled", func() {
			installation.CertificateManagement = &operatorv1.CertificateManagement{}
			expectedResources := []resourceTestObj{
				{relasticsearch.PublicCertSecret, rmeta.OperatorNamespace(), &corev1.Secret{}, nil},
				{render.TigeraElasticsearchCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{relasticsearch.PublicCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{render.KibanaInternalCertSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersUserSecret, rmeta.OperatorNamespace(), &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersVerificationUserSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{kubecontrollers.ElasticsearchKubeControllersSecureUserSecret, render.ElasticsearchNamespace, &corev1.Secret{}, nil},
				{ServiceName, render.ElasticsearchNamespace, &corev1.Service{}, nil},
				{RoleName, render.ElasticsearchNamespace, &rbacv1.Role{}, nil},
				{RoleName, render.ElasticsearchNamespace, &rbacv1.RoleBinding{}, nil},
				{ServiceAccountName, render.ElasticsearchNamespace, &corev1.ServiceAccount{}, nil},
				{DeploymentName, render.ElasticsearchNamespace, &appsv1.Deployment{}, nil},
				{RoleName + ":csr-creator", "", &rbacv1.ClusterRoleBinding{}, nil},
			}

			component := EsGateway(&Config{
				installation,
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: "tigera-pull-secret"}},
				},
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: render.TigeraElasticsearchCertSecret, Namespace: rmeta.OperatorNamespace()}},
					{ObjectMeta: metav1.ObjectMeta{Name: relasticsearch.PublicCertSecret, Namespace: rmeta.OperatorNamespace()}},
				},
				[]*corev1.Secret{
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersUserSecret, Namespace: rmeta.OperatorNamespace()}},
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersVerificationUserSecret, Namespace: render.ElasticsearchNamespace}},
					{ObjectMeta: metav1.ObjectMeta{Name: kubecontrollers.ElasticsearchKubeControllersSecureUserSecret, Namespace: render.ElasticsearchNamespace}},
				},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: render.KibanaInternalCertSecret, Namespace: rmeta.OperatorNamespace()}},
				&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: relasticsearch.InternalCertSecret, Namespace: render.ElasticsearchNamespace}},
				clusterDomain, "elastic",
			})

			createResources, _ := component.Objects()
			compareResources(createResources, expectedResources)
		})

	})
})

func compareResources(resources []client.Object, expectedResources []resourceTestObj) {
	Expect(len(resources)).To(Equal(len(expectedResources)))
	for i, expectedResource := range expectedResources {
		resource := resources[i]
		actualName := resource.(metav1.ObjectMetaAccessor).GetObjectMeta().GetName()
		actualNS := resource.(metav1.ObjectMetaAccessor).GetObjectMeta().GetNamespace()

		Expect(actualName).To(Equal(expectedResource.name), fmt.Sprintf("Rendered resource has wrong name (position %d, name %s, namespace %s)", i, actualName, actualNS))
		Expect(actualNS).To(Equal(expectedResource.ns), fmt.Sprintf("Rendered resource has wrong namespace (position %d, name %s, namespace %s)", i, actualName, actualNS))
		Expect(resource).Should(BeAssignableToTypeOf(expectedResource.typ))
		if expectedResource.f != nil {
			expectedResource.f(resource)
		}
	}
}
