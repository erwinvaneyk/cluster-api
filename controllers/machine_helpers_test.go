/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_getActiveMachinesInCluster(t *testing.T) {
	ns1Cluster1 := clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ns1cluster1",
			Namespace: "test-ns-1",
			Labels: map[string]string{
				clusterv1.MachineClusterLabelName: "test-cluster-1",
			},
		},
	}
	ns1Cluster2 := clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ns1cluster2",
			Namespace: "test-ns-1",
			Labels: map[string]string{
				clusterv1.MachineClusterLabelName: "test-cluster-2",
			},
		},
	}
	time := metav1.Now()
	ns1Cluster1Deleted := clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ns1cluster1deleted",
			Namespace: "test-ns-1",
			Labels: map[string]string{
				clusterv1.MachineClusterLabelName: "test-cluster-2",
			},
			DeletionTimestamp: &time,
		},
	}
	ns2Cluster2 := clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ns2cluster2",
			Namespace: "test-ns-2",
			Labels: map[string]string{
				clusterv1.MachineClusterLabelName: "test-cluster-2",
			},
		},
	}

	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*clusterv1.Machine
		wantErr bool
	}{
		{
			name: "ns1 cluster1",
			args: args{
				namespace: "test-ns-1",
				name:      "test-cluster-1",
			},
			want:    []*clusterv1.Machine{&ns1Cluster1},
			wantErr: false,
		},
		{
			name: "ns2 cluster2",
			args: args{
				namespace: "test-ns-2",
				name:      "test-cluster-2",
			},
			want:    []*clusterv1.Machine{&ns2Cluster2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fake.NewFakeClient(&ns1Cluster1, &ns1Cluster2, &ns1Cluster1Deleted, &ns2Cluster2)
			got, err := getActiveMachinesInCluster(context.TODO(), c, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getActiveMachinesInCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getActiveMachinesInCluster() got = %v, want %v", got, tt.want)
			}
		})
	}
}
