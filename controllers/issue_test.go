package controllers

import (
	"context"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeAnnotationValue(t *testing.T) {
	const key = "field.cattle.io/publicEndpoints"
	settings := provideSettings(t, predefinedConfig)
	settings.ingress.Annotations[key] = "1.1.1.1"
	settings.client.Update(context.TODO(), settings.ingress)
	settings.reconciler.createGSLBFromIngress(settings.client, strategyAnnotation, depresolver.RoundRobinStrategy, settings.ingress, depresolver.RoundRobinStrategy)
	_ = settings.reconciler.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	require.Equal(t, len(settings.gslb.Annotations), 2)
	require.Equal(t, settings.gslb.Annotations[key], "1.1.1.1")

	settings.ingress.Annotations[key] = "1.1.1.2"
	settings.client.Update(context.Background(), settings.ingress)
	settings.reconciler.createGSLBFromIngress(settings.client, strategyAnnotation, depresolver.RoundRobinStrategy, settings.ingress, depresolver.RoundRobinStrategy)
	_ = settings.reconciler.Get(context.TODO(), settings.request.NamespacedName, settings.gslb)
	require.Equal(t, len(settings.gslb.Annotations), 2)
	require.Equal(t, settings.gslb.Annotations[key], "1.1.1.2")
}

func TestChangeAnnotationValueWithReconcile(t *testing.T) {
	const key = "field.cattle.io/publicEndpoints"
	settings := provideSettings(t, predefinedConfig)
	settings.ingress.Annotations[key] = "1.1.1.1"
	settings.client.Update(context.TODO(), settings.ingress)
	settings.reconciler.createGSLBFromIngress(settings.client, strategyAnnotation, depresolver.RoundRobinStrategy, settings.ingress, depresolver.RoundRobinStrategy)
	reconcileAndUpdateGslb(t, settings)
	require.Equal(t, len(settings.gslb.Annotations), 2)
	require.Equal(t, settings.gslb.Annotations[key], "1.1.1.1")

	settings.ingress.Annotations[key] = "1.1.1.2"
	settings.client.Update(context.Background(), settings.ingress)
	settings.reconciler.createGSLBFromIngress(settings.client, strategyAnnotation, depresolver.RoundRobinStrategy, settings.ingress, depresolver.RoundRobinStrategy)
	reconcileAndUpdateGslb(t, settings)
	require.Equal(t, len(settings.gslb.Annotations), 2)
	require.Equal(t, settings.gslb.Annotations[key], "1.1.1.2")
}

func TestChangeIngressHostName(t *testing.T) {
	const host = "found.cloud.example.com"
	settings := provideSettings(t, predefinedConfig)
	settings.ingress.Spec.Rules[0].Host = host
	settings.client.Update(context.TODO(), settings.ingress)
	settings.reconciler.createGSLBFromIngress(settings.client, strategyAnnotation, depresolver.RoundRobinStrategy, settings.ingress, depresolver.RoundRobinStrategy)
	reconcileAndUpdateGslb(t, settings)
	require.Equal(t, settings.gslb.Spec.Ingress.Rules[0].Host, settings.ingress.Spec.Rules[0].Host)
	require.Equal(t, settings.gslb.Spec.Ingress.Rules[1].Host, settings.ingress.Spec.Rules[1].Host)
	require.Equal(t, settings.gslb.Spec.Ingress.Rules[2].Host, settings.ingress.Spec.Rules[2].Host)
	require.Equal(t, 3, len(settings.gslb.Spec.Ingress.Rules))
	require.Equal(t, 3, len(settings.ingress.Spec.Rules))
}