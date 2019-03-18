/*
This package contains Dynamic Kubernetes Host Path PV provisioner
using external storage provisioner library:
https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner

Host Path volume mount is one of the oldest supported option for Kubernetes
Stateful workloads and makes a great choice in the following scenarios:
- The Stateful Workload need to persist data across node reboots
  and/or pod restarts.
- The Stateful Workload can take care of replicating the data across
  nodes to handle cases like a complete node (and/or its storage) failure.
- For long running Stateful Workloads, the Backup/Recovery is provided
  by Operators/tools that can make use the Workload mounts and do not
  require the capabilities to be available in the underlying storage. Or
  if the host paths are created on external storage like EBS/GPD, administrator
  have tools that can periodically take snapshots/backups.

However, the use of host path by application developers has been discouraged
due to potential security vulnerabilities of application developer gaining
access (inadvertently) to other directories on the host and causing irrepairable
damage. We are seeing that enterprise now setup Pod Security Policy that restricts the
application developers from using host paths or at least restrict to only trusted
namespaces. For more, check out:
https://kubernetes.io/docs/concepts/policy/pod-security-policy/#volumes-and-file-systems

The Local PVs provide an viable alternative to using the host paths,
while addressing the security risk of exposing the host paths to the
application developers. For example application developer can now create
PVCs - that request for host path based Local PV.

Host Path Provisioner can dynamically satisfy the PVC request
by providing a PV object. The PV objects are not accessible by the
developers and hence they are shielded from the actual paths used
on the host.

While the primary motivation to add Host Path PVs was to address the
security vulnerabilities, host paths are probably the only viable options
in many other small to medium sized Kubernetes Clusters with limited amount
of resources(disks/storage or CPU/RAM to run storage engines) available per node.

Host Path (like Local PVs) have become ready to be used with Kubernetes 1.12
storage features like Local Volume Topology and the fact that Pod scheduling
can take into account the PV NodeAffinity parameters.

Host Path PV vs Local PV:
- Multiple Host Path PVs can be provisioned on a single underlying disk. Local PVs are designed to use up the entire disk.
- Host Path PV only supports filesystem mounts. Local PVs can also support block mounts.

Host Path PV vs native hostPath:
- Host Path PV supports volume topology constraints, allowing the Pods using the Host Path PV to 
  be pinned to the node. (Can also be achieved by using a node selector on the pods, but can't be 
  scaled. Also application developers need not worry about selecting the node)
- Host Path PV supports enforcing checks on the paths used. (Can also be achieved via PSP, but 
  enforcing via StorageClass makes it easy to put all storage policies at one place.)
- Host Path PV support backup/restore of data using tools like Heptio Velero
- Host Path PV supports enforcing the Capacity limits. Host Path PV can ve converted to 
  Read-Only when the Capacity limit is reached.
- Host Path PV supports selecting a node with the available capacity. Can limit the number of 
  Host Path PVs provisioned based on the capacity available on the node.

How it works:
Host Path Provisioner uses the StorageClass `waitForConsumer` option, so it
- waits for the workload to be scheduled to a suitable node.
- on the selected node, a PV object is created with a unique host path

The exact location of the host path can be specified using the following options:
- BasePath: The Host Path volumes will be created under `/var/openebs`.
  This default path can be changed by passing the "OPENEBS_IO_BASE_PATH" ENV
  variable to the Host Path Provisioner Pod.
  It is also possible to specify a different location using the CAS Policy
  `BasePath` in the StorageClass or the PVC.
  The CAS Policy can either be provided by the Administrator or could also be injected
  via the Policy Agents.
- RelativePath: CAS Policy can help with cases where an application have to be loaded using
  a previously shared or known path - either for recovery or maintenance purposes.

The location of the host paths can be controlled by the cluster administrator
using the configuration options mentioned above and can be any of the following:
- OS Disk  - possibly a folder dedicated to saving data on each node.
- Additional Disks - mounted as ext4 or any other FS of choice on the hosts
- External Storage - mounted as ext4 or any other FS of choice on the hosts.

Limitations:
- It is not possible to control the nodes where host paths can be created.
  This can be only be done by controlling where the stateful workload are scheduled.
- It is not possible to Limit the PVs created at a given path based on the available storage.
  The initial implementation doesn't honor the capacity limits of the PVC.
- Noisy neighbor effects can be seen if multiple workloads share the same underlying disks.


Inspiration:
The implementation has been influenced by the prior work done by the Kubernetes community,
specifically the following:
- https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/tree/master/examples/hostpath-provisioner
- https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner
- https://github.com/rancher/local-path-provisioner

Future Work:
Host Paths provide an easy (almost no hassle) way of running workloads. The build of the
storage Day 0 (setting up of the mount paths) and Day 2 (expansion and migration) operations
can be handled with projects like NDM or others. Some of the future ideas to be implement:
- Ability to use a git backed host path, so data can be backed up to a github/gitlab
- Ability to use host paths that are managed by NDM - that monitors for usage, helps
  with expanding the storage of a given host path. For example - use LVM or ZFS to
  create a host mount with attached disks. Where additional disks can be added or failed
  disks replaced without impacting the workloads running on their host paths.
- Integrate with Projects like Valero or Kasten that can handle backup and restor of data
  stored on the Host Path PVs attached to a workload.
- Provide tools that can help with recovering from situations where PVs are tied to nodes
  that can never recover. For example, a Stateful Workload can be associated with a
  Host Path PV on Node-z. Say Node-z  becomes inaccassible for reasons beyond the control
  like - a site/zone/rack disaster or the disks went up in flames. The PV will still be
  having the Node Affinity to Node-z, which will make the Workload to get stuck in pending
  state.
- Move towards using a CSI based Host Path provisioner. Some of the features required to
  use the Host Path PVs like the Volume Topology are not yet available in the CSI. As the
  CSI driver stabilizes, this can be moved into CSI.

*/
package app
