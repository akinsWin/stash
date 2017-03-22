package test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/appscode/restik/pkg/controller"
	"github.com/stretchr/testify/assert"
)

func TestBackups(t *testing.T) {
	randStr := RandStringRunes(10)
	var backupRC = "backup-test-replicationcontroller-" + randStr
	var backupReplicaset = "backup-test-replicaset-" + randStr
	var backupDeployment = "backup-test-deployment-" + randStr
	var backupDaemonset = "backup-test-daemonset-" + randStr
	var backupStatefulset = "backup-test-statefulset-" + randStr
	var repoSecret = "restik-test-secret-" + randStr
	var rs = "restik-test-replicaset-" + randStr
	var deployment = "restik-test-deployment-" + randStr
	var rc = "restik-test-rc-" + randStr
	var statefulset = "restik-test-statefulset-" + randStr
	var svc = "restik-test-svc-" + randStr
	var daemonset = "restik-test-daemonset-" + randStr
	log.Println("###############==Running e2e tests for Restik==#############")
	watcher, err := runController()
	if !assert.Nil(t, err) {
		return
	}
	namespace := "test-" + randStr
	log.Println("\n\nCreating namespace -->", namespace)
	err = createTestNamespace(watcher, namespace)
	if !assert.Nil(t, err) {
		return
	}
	defer deleteTestNamespace(watcher, namespace)
	log.Println("Creating secret with password -->", repoSecret)
	err = createSecret(watcher, repoSecret)
	if !assert.Nil(t, err) {
		return
	}
	defer deleteSecret(watcher, repoSecret)

	log.Println("\n***********************************************************\nCreating Daemonset -->", daemonset)
	err = createDaemonsets(watcher, daemonset, backupDaemonset)
	if !assert.Nil(t, err) {
		return
	}
	time.Sleep(time.Second * 10)
	log.Printf("Starting backup(%s) for Daemonset...\n", backupDaemonset)
	err = createBackup(watcher, backupDaemonset, repoSecret)
	if !assert.Nil(t, err) {
		return
	}
	err = checkEventForBackup(watcher, backupDaemonset+"-1")
	if !assert.Nil(t, err) {
		return
	}
	log.Println("Removing backup for Daemonset")
	err = deleteBackup(watcher, backupDaemonset)
	if !assert.Nil(t, err) {
		return
	}
	err = checkContainerAfterBackupDelete(watcher, daemonset, controller.DaemonSet)
	if !assert.Nil(t, err) {
		return
	}
	log.Println("SUCCESS: Daemonset Backup")
	deleteDaemonset(watcher, daemonset)
	log.Println("\n************************************************************\nCreating ReplicationController -->", rc)
	err = createReplicationController(watcher, rc, backupRC)
	if !assert.Nil(t, err) {
		return
	}
	time.Sleep(time.Second * 10)
	log.Printf("Starting backup(%s) for ReplicationController...\n", backupRC)
	err = createBackup(watcher, backupRC, repoSecret)
	if !assert.Nil(t, err) {
		return
	}

	err = checkEventForBackup(watcher, backupRC+"-1")
	if !assert.Nil(t, err) {
		return
	}
	log.Println("Removing backup for ReplicationController")
	err = deleteBackup(watcher, backupRC)
	if !assert.Nil(t, err) {
		return
	}
	err = checkContainerAfterBackupDelete(watcher, rc, controller.ReplicationController)
	if !assert.Nil(t, err) {
		return
	}
	log.Println("SUCCESS: ReplicationController Backup")
	deleteReplicationController(watcher, rc)

	log.Println("\n***********************************************************\nCreating Replicaset -->", rs)
	err = createReplicaset(watcher, rs, backupReplicaset)
	if !assert.Nil(t, err) {
		return
	}
	time.Sleep(time.Second * 10)
	log.Printf("Starting backup(%s) for Replicaset...\n", backupReplicaset)
	err = createBackup(watcher, backupReplicaset, repoSecret)
	if !assert.Nil(t, err) {
		return
	}
	err = checkEventForBackup(watcher, backupReplicaset+"-1")
	if !assert.Nil(t, err) {
		return
	}
	log.Println("Removing backup for Replicaset")
	err = deleteBackup(watcher, backupReplicaset)
	if !assert.Nil(t, err) {
		return
	}
	err = checkContainerAfterBackupDelete(watcher, rs, controller.ReplicaSet)
	if !assert.Nil(t, err) {
		return
	}
	log.Println("SUCCESS : Replicaset Backup")
	deleteReplicaset(watcher, rs)

	log.Println("\n***********************************************************\nCreating Deployment -->", deployment)
	err = createDeployment(watcher, deployment, backupDeployment)
	if !assert.Nil(t, err) {
		return
	}
	time.Sleep(time.Second * 10)
	defer deleteDeployment(watcher, deployment)
	log.Printf("Starting backup(%s) for deployment...\n", backupDeployment)
	err = createBackup(watcher, backupDeployment, repoSecret)
	if !assert.Nil(t, err) {
		return
	}
	err = checkEventForBackup(watcher, backupDeployment+"-1")
	if !assert.Nil(t, err) {
		return
	}
	log.Println("Removing backup for Deployment")
	err = deleteBackup(watcher, backupDeployment)
	if !assert.Nil(t, err) {
		return
	}
	err = checkContainerAfterBackupDelete(watcher, deployment, controller.Deployment)
	if !assert.Nil(t, err) {
		return
	}
	log.Println("SUCCESS: Deployment Backup")

	log.Println("\n***********************************************************\nCreating Statefulset with Backup -->", statefulset)
	err = createService(watcher, svc)
	if !assert.Nil(t, err) {
		return
	}
	err = createStatefulSet(watcher, statefulset, backupStatefulset, svc)
	if !assert.Nil(t, err) {
		return
	}
	time.Sleep(time.Second * 10)
	err = createBackup(watcher, backupStatefulset, repoSecret)
	if !assert.Nil(t, err) {
		return
	}
	defer deleteBackup(watcher, backupStatefulset)

	err = checkEventForBackup(watcher, backupStatefulset+"-1")
	if !assert.Nil(t, err) {
		return
	}
	log.Println("SUCCESS: Statefulset Backup")
	deleteStatefulset(watcher, statefulset)
	deleteService(watcher, svc)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}