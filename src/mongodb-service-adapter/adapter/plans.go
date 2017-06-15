package adapter

var plans = map[string]string{
	"standalone": `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation",
        "downloadBaseWindows": "C:\\mongodb-mms-automation"
    },
    "mongoDbVersions": [{
        "builds": [
                {
                    "bits": 64,
                    "flavor": "",
                    "gitVersion": "4249c1d2b5999ebbf1fdf3bc0e0e3b3ff5c0aaf2",
                    "maxOsVersion": "",
                    "minOsVersion": "",
                    "modules": [],
                    "platform": "osx",
                    "url": "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon-3.2.7.tgz",
                    "win2008plus": false,
                    "winVCRedistDll": "",
                    "winVCRedistOptions": [],
                    "winVCRedistUrl": "",
                    "winVCRedistVersion": ""
                }
            ],
        "name": "{{version}}"
    }],
    "backupVersions": [{
        "hostname": "{{nodes.[0]}}",
        "logPath": "/var/vcap/sys/log/mongod_node/backup-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],

    "monitoringVersions": [{
        "hostname": "{{nodes.[0]}}",
        "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],
    "processes": [{
        "args2_6": {
            "net": {
                "port": 28000
            },
            "storage": {
                "dbPath": "/var/vcap/store/mongodb-data"
            },
            "systemLog": {
                "destination": "file",
                "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
            }
        },
        "hostname": "{{nodes.[0]}}",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        },
        "name": "{{nodes.[0]}}",
        "processType": "mongod",
        "version": "{{version}}",
        "authSchemaVersion": 5
    }],
    "replicaSets": [],
    "roles": [],
    "sharding": [],

    "auth": {
        "autoUser": "mms-automation",
        "autoPwd": "{{ password }}",
        "deploymentAuthMechanisms": [
            "SCRAM-SHA-1"
        ],
        "key": "{{ key }}",
        "keyfile": "/var/vcap/jobs/mongod_node/config/mongo_om.key",
        "disabled": false,
        "usersDeleted": [],
        "usersWanted": [
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterMonitor"
                    }
                ],
                "user": "mms-monitoring-agent",
                "initPwd": "{{ password }}"
            },
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterAdmin"
                    },
                    {
                        "db": "admin",
                        "role": "readAnyDatabase"
                    },
                    {
                        "db": "admin",
                        "role": "userAdminAnyDatabase"
                    },
                    {
                        "db": "local",
                        "role": "readWrite"
                    },
                    {
                        "db": "admin",
                        "role": "readWrite"
                    }
                ],
                "user": "mms-backup-agent",
                "initPwd": "{{ password }}"
            },
            {
               "db": "admin" ,
               "user": "admin" ,
               "roles": [
                 {
                     "db": "admin",
                     "role": "clusterAdmin"
                 },
                 {
                     "db": "admin",
                     "role": "readAnyDatabase"
                 },
                 {
                     "db": "admin",
                     "role": "userAdminAnyDatabase"
                 },
                 {
                     "db": "local",
                     "role": "readWrite"
                 },
                 {
                     "db": "admin",
                     "role": "readWrite"
                 }
               ],
               "initPwd": "{{ admin_password }}"
            }
        ],
        "autoAuthMechanism": "SCRAM-SHA-1"
    }
}`,

	"sharded_set": `{
   "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation",
        "downloadBaseWindows": "C:\\mongodb-mms-automation"
    },
    "mongoDbVersions": [{
        "builds": [
                {
                    "bits": 64,
                    "flavor": "",
                    "gitVersion": "4249c1d2b5999ebbf1fdf3bc0e0e3b3ff5c0aaf2",
                    "maxOsVersion": "",
                    "minOsVersion": "",
                    "modules": [],
                    "platform": "osx",
                    "url": "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon-3.2.7.tgz",
                    "win2008plus": false,
                    "winVCRedistDll": "",
                    "winVCRedistOptions": [],
                    "winVCRedistUrl": "",
                    "winVCRedistVersion": ""
                }
            ],
        "name": "{{version}}"
    }],
    "backupVersions": [
    ],

    "monitoringVersions": [
    {
        "hostname": "{{nodes.[0]}}",
        "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }
    ],
    "processes": [{{#each nodes}}
      {
        "args2_6": {
            "net": {
                "port": 28000
            },
            {{#if (isInShard @index)}}
            "replication": {
                "replSetName": "shard_{{div @index 3 }}"
            },
            {{/if}}
            {{#if (isConfig @index)}}
            "sharding": {
                "clusterRole": "configsvr"
            },
            {{/if}}
            {{#if (hasStorage @index)}}
            "storage": {
                "dbPath": "/var/vcap/store/mongodb-data"
            },
            {{/if}}
            "systemLog": {
                "destination": "file",
                "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
            }
        },
        "hostname": "{{this}}",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        },
        "name": "{{this}}",
        "processType": "{{processType @index}}",
        {{#if (hasShardedCluster @index)}}
        "cluster": "sharded-cluster",
        {{/if}}
        "version": "{{version}}",
        "authSchemaVersion": 5
    }{{#if @last}}{{else}},{{/if}}
    {{/each}}
    ],
    "replicaSets": [{{#each partitionedNodes}}

            {
                "_id": "shard_{{@index}}",
                "members": [{{#each this}}
                    {
                        "_id": {{@index}},
                        "arbiterOnly": false,
                        "hidden": false,
                        "host": "{{this}}",
                        "priority": 1,
                        "slaveDelay": 0,
                        "votes": 1
                    }{{#if @last}}{{else}},{{/if}}
                    {{/each}}
                ]
            },
            {{/each}}
            {
                "_id": "shard_3",
                "members": [
                    {
                        "_id": 0,
                        "arbiterOnly": false,
                        "hidden": false,
                        "host": "{{nodes.[9]}}",
                        "priority": 1,
                        "slaveDelay": 0,
                        "votes": 1
                    },
                    {
                        "_id": 1,
                        "arbiterOnly": false,
                        "hidden": false,
                        "host": "{{nodes.[10]}}",
                        "priority": 1,
                        "slaveDelay": 0,
                        "votes": 1
                    },
                    {
                        "_id": 2,
                        "arbiterOnly": false,
                        "hidden": false,
                        "host": "{{nodes.[11]}}",
                        "priority": 1,
                        "slaveDelay": 0,
                        "votes": 1
                    }
                ]
            }
    ],
    "sharding": [
        {
                "shards": [
                  {{#each partitionedNodes}}
                    {
                        "tags": [],
                        "_id": "shard_{{@index}}",
                        "rs": "shard_{{@index}}"
                    }{{#if @last}}{{else}},{{/if}}
                    {{/each}}
                ],
                "name": "sharded-cluster",
                "configServer": [],
                "configServerReplica": "shard_3",
                "collections": []
            }
    ],

    "auth":{
       "disabled":false,
       "autoPwd": "{{ password }}",
       "autoUser":"mms-automation",
       "deploymentAuthMechanisms": [
           "MONGODB-CR"
       ],
       "key":"{{ key }}",
       "keyfile":"/var/vcap/jobs/mongod_node/config/mongo_om.key",
       "usersWanted":[
          {
             "db":"admin",
             "initPwd":"{{ password }}",
             "roles":[
                {
                   "db":"admin",
                   "role":"clusterMonitor"
                }
             ],
             "user":"mms-monitoring-agent"
          },
          {
             "db":"admin",
             "initPwd":"{{ password }}",
             "roles":[
                {
                   "db":"admin",
                   "role":"clusterAdmin"
                },
                {
                   "db":"admin",
                   "role":"readAnyDatabase"
                },
                {
                   "db":"admin",
                   "role":"userAdminAnyDatabase"
                },
                {
                   "db":"local",
                   "role":"readWrite"
                },
                {
                   "db":"admin",
                   "role":"readWrite"
                }
             ],
             "user":"mms-backup-agent"
          }
       ],
       "usersDeleted":[],
       "autoAuthMechanism": "MONGODB-CR"
    }
}`,

	"single_replica_set": `{
    "options": {
        "downloadBase": "/var/lib/mongodb-mms-automation",
        "downloadBaseWindows": "C:\\mongodb-mms-automation"
    },
    "mongoDbVersions": [{
        "builds": [
                {
                    "bits": 64,
                    "flavor": "",
                    "gitVersion": "4249c1d2b5999ebbf1fdf3bc0e0e3b3ff5c0aaf2",
                    "maxOsVersion": "",
                    "minOsVersion": "",
                    "modules": [],
                    "platform": "osx",
                    "url": "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon-3.2.7.tgz",
                    "win2008plus": false,
                    "winVCRedistDll": "",
                    "winVCRedistOptions": [],
                    "winVCRedistUrl": "",
                    "winVCRedistVersion": ""
                }
            ],
        "name": "{{version}}"
    }],
    "backupVersions": [{
        "hostname": "{{nodes.[0]}}",
        "logPath": "/var/vcap/sys/log/mongod_node/backup-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],

    "monitoringVersions": [{
        "hostname": "{{nodes.[0]}}",
        "logPath": "/var/vcap/sys/log/mongod_node/monitoring-agent.log",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        }
    }],
    "processes": [{{#each nodes}}
      {
        "args2_6": {
            "net": {
                "port": 28000
            },
            "replication": {
                "replSetName": "pcf_repl"
            },
            "storage": {
                "dbPath": "/var/vcap/store/mongodb-data"
            },
            "systemLog": {
                "destination": "file",
                "path": "/var/vcap/sys/log/mongod_node/mongodb.log"
            }
        },
        "hostname": "{{this}}",
        "logRotate": {
            "sizeThresholdMB": 1000,
            "timeThresholdHrs": 24
        },
        "name": "{{this}}",
        "processType": "mongod",
        "version": "{{version}}",
        "authSchemaVersion": 5
    }{{#if @last}}{{else}},{{/if}}
    {{/each}}
  ],
    "replicaSets": [{
        "_id": "pcf_repl",
        "members": [
          {{#each nodes}}
          {
            "_id": {{@index}},
            "host": "{{this}}"
          {{#if @last}},"arbiterOnly": true,"priority": 0}
          {{else}}},{{/if}}{{/each}}
        ]
    }],
    "roles": [],
    "sharding": [],

    "auth": {
        "autoUser": "mms-automation",
        "autoPwd": "{{ password }}",
        "deploymentAuthMechanisms": [
            "SCRAM-SHA-1"
        ],
        "key": "{{ key }}",
        "keyfile": "/var/vcap/jobs/mongod_node/config/mongo_om.key",
        "disabled": false,
        "usersDeleted": [],
        "usersWanted": [
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterMonitor"
                    }
                ],
                "user": "mms-monitoring-agent",
                "initPwd": "{{ password }}"
            },
            {
                "db": "admin",
                "roles": [
                    {
                        "db": "admin",
                        "role": "clusterAdmin"
                    },
                    {
                        "db": "admin",
                        "role": "readAnyDatabase"
                    },
                    {
                        "db": "admin",
                        "role": "userAdminAnyDatabase"
                    },
                    {
                        "db": "local",
                        "role": "readWrite"
                    },
                    {
                        "db": "admin",
                        "role": "readWrite"
                    }
                ],
                "user": "mms-backup-agent",
                "initPwd": "{{ password }}"
            },
            {
               "db": "admin" ,
               "user": "admin" ,
               "roles": [
                 {
                     "db": "admin",
                     "role": "clusterAdmin"
                 },
                 {
                     "db": "admin",
                     "role": "readAnyDatabase"
                 },
                 {
                     "db": "admin",
                     "role": "userAdminAnyDatabase"
                 },
                 {
                     "db": "local",
                     "role": "readWrite"
                 },
                 {
                     "db": "admin",
                     "role": "readWrite"
                 }
               ],
               "initPwd": "{{ admin_password }}"
            }
        ],
        "autoAuthMechanism": "SCRAM-SHA-1"
    }
}`,
}
