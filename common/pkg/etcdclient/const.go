package etcdclient

const (
	keyEtcdProfile = "/crony/"

	//node节点
	//key /crony/node/<node_uuid>
	KeyEtcdNodeProfile = keyEtcdProfile + "node/"
	KeyEtcdNode        = KeyEtcdNodeProfile + "%s"

	//key  /crony/proc/<node_uuid>/<job_id>/<pid>
	KeyEtcdProcProfile = keyEtcdProfile + "proc/%s"
	KeyEtcdProc        = KeyEtcdProcProfile + "%d/%d"

	//key /crony/job/<node_uuid>/<job_id>
	KeyEtcdJobProfile = keyEtcdProfile + "job/%s/"
	KeyEtcdJob        = KeyEtcdJobProfile + "%d"

	// key /crony/once/<jobID>
	KeyEtcdOnceProfile = keyEtcdProfile + "once/"
	KeyEtcdOnce        = KeyEtcdOnceProfile + "%d"

	KeyEtcdLockProfile = keyEtcdProfile + "lock/"
	KeyEtcdLock        = KeyEtcdLockProfile + "%s"

	// key /crony/system/<node_uuid>
	KeyEtcdSystemProfile = keyEtcdProfile + "system/"
	KeyEtcdSystemSwitch  = KeyEtcdSystemProfile + "switch/" + "%s"
	KeyEtcdSystemGet     = KeyEtcdSystemProfile + "get/" + "%s"
)
