class Resolver {
  resolveVariable;
  resolveConfigurationProperty;

  constructor(resolveVariable, resolveConfigurationProperty) {
    this.resolveVariable = resolveVariable;
    this.resolveConfigurationProperty = resolveConfigurationProperty;
  }

  async variable(variableStr) {
    return await this.resolveVariable(variableStr);
  }

  async configProperty(...keys) {
    return await this.resolveConfigurationProperty(keys);
  }

  async resolve(namespaces) {
    let variables = {};
    for (const ns in namespaces) {
      variables[ns] = {};
      for (const variable of namespaces[ns]) {
        const value = await this.resolveVariable(`${ns}:${variable}`);
        let current = variables[ns];
        const parts = variable.split(".");
        parts.forEach((part, index) => {
          if (index < parts.length - 1) {
            current[part] = {};
            current = current[part];
          } else {
            current[part] = value;
          }
        });
      }
    }
    return variables;
  }

  async constructSLS() {
    let sls = new SLS();
    sls.variables = await this.resolve({
      self: ["service"],
      sls: ["stage"],
    });
    return sls;
  }
}

class SLS {
  variables;

  service() {
    return this.variables.self.service;
  }

  stage() {
    return this.variables.sls.stage;
  }

  fullName(name) {
    return `${this.service()}-${this.stage()}-${name}`;
  }

  getAtt(logicalName, attributeName) {
    return { "Fn::GetAtt": [logicalName, attributeName] };
  }

  join(...parts) {
    return { "Fn::Join": ["", parts] };
  }

  getArn(logicalName) {
    return this.getAtt(logicalName, "Arn");
  }

  parseLogicalPath(logicalPath) {
    return /([^\/]+)(\/.+)?/.exec(logicalPath).slice(1, 3);
  }

  Resource(logicalPath) {
    return new Resource(this, [this.parseLogicalPath(logicalPath)]);
  }

  Resources(...logicalPaths) {
    return new Resource(
      this,
      logicalPaths.map(this.parseLogicalPath.bind(this))
    );
  }
}

const Actions = (() => {
  const actions = {};

  actions.merge = (...acts) => {
    return [...new Set(acts.flat())];
  };

  actions.DynamoDBRead = [
    "dynamodb:GetItem",
    "dynamodb:BatchGetItem",
    "dynamodb:Scan",
    "dynamodb:Query",
    "dynamodb:ConditionCheckItem",
  ];

  actions.DynamoDBCrud = actions.merge(actions.DynamoDBRead, [
    "dynamodb:DeleteItem",
    "dynamodb:PutItem",
    "dynamodb:UpdateItem",
    "dynamodb:BatchWriteItem",
    "dynamodb:DescribeTable",
  ]);

  actions.S3Read = [
    "s3:GetObject",
    "s3:ListBucket",
    "s3:GetBucketLocation",
    "s3:GetObjectVersion",
    "s3:GetLifecycleConfiguration",
  ];

  actions.S3Crud = actions.merge(actions.S3Read, [
    "s3:PutObject",
    "s3:PutObjectAcl",
    "s3:PutLifecycleConfiguration",
    "s3:DeleteObject",
  ]);

  return actions;
})();

class Resource {
  arn;
  sls;

  constructor(sls, resources) {
    this.sls = sls;
    if (resources.length > 1) {
      this.arn = resources.map(this.qualify.bind(this));
    } else {
      this.arn = this.qualify(resources[0]);
    }
  }

  qualify(resource) {
    let [logicalName, path] = resource;
    let arn = this.sls.getArn(logicalName);
    if (path !== undefined) {
      arn = this.sls.join(arn, path);
    }
    return arn;
  }

  Policy(...actions) {
    return new ResourcePolicy(this, Actions.merge(...actions));
  }
}

class ResourcePolicy {
  resource;

  Effect;
  Action;
  Resource;

  constructor(resource, actions) {
    Object.defineProperty(this, "resource", {
      enumerable: false,
      value: resource,
    });
    Object.freeze(this.resource);

    this.Effect = "Allow";
    this.Action = actions;
    this.Resource = this.resource.arn;
  }

  Policy(...actions) {
    this.Action = Actions.merge(...this.Action, ...actions);
    return this;
  }
}

const HttpApi = (() => {
  let httpApi = {};
  const call = (method) => (path) => {
    return {
      httpApi: {
        method,
        path,
      },
    };
  };

  httpApi.Post = call("Post");
  httpApi.Get = call("Get");
  httpApi.Delete = call("Delete");
  httpApi.Put = call("Put");
  return httpApi;
})();

module.exports = async ({ resolveVariable, resolveConfigurationProperty }) => {
  const resolver = new Resolver(resolveVariable, resolveConfigurationProperty);
  const sls = await resolver.constructSLS();

  const handler = (name) =>
    `github.com/renbou/dontstress/lambda-api/handlers/${name}`;

  return {
    addTaskDataFunction: {
      handler: handler("add-task-data"),
      iamRoleStatements: [
        sls.Resource("adminsTable").Policy(Actions.DynamoDBRead),
        sls.Resource("tasksTable").Policy(Actions.DynamoDBCrud),
        sls.Resource("filesTable").Policy(Actions.DynamoDBCrud),
        sls.Resources("filesBucket", "filesBucket/*").Policy(Actions.S3Crud),
      ],
      events: [HttpApi.Post("/lab/{labid}/task/{taskid}")],
    },

    bebraFunction: {
      handler: handler("bebra"),
      events: [HttpApi.Get("/bebra"), HttpApi.Post("/bebra")],
    },

    createLabFunction: {
      handler: handler("create-lab"),
      iamRoleStatements: [
        sls.Resource("adminsTable").Policy(Actions.DynamoDBRead),
        sls.Resource("labsTable").Policy(Actions.DynamoDBCrud),
      ],
      events: [HttpApi.Post("/labs")],
    },

    createTaskFunction: {
      handler: handler("create-lab-task"),
      iamRoleStatements: [
        sls.Resource("adminsTable").Policy(Actions.DynamoDBRead),
        sls.Resource("tasksTable").Policy(Actions.DynamoDBCrud),
        sls.Resource("labsTable").Policy(Actions.DynamoDBRead),
      ],
      events: [HttpApi.Post("/lab/{labid}/tasks")],
    },

    deleteLabFunction: {
      handler: handler("delete-lab"),
      iamRoleStatements: [
        sls.Resource("adminsTable").Policy(Actions.DynamoDBRead),
        sls.Resource("labsTable").Policy(Actions.DynamoDBCrud),
      ],
      events: [HttpApi.Delete("/lab/{labid}")],
    },

    deleteTaskFunction: {
      handler: handler("delete-task"),
      iamRoleStatements: [
        sls.Resource("adminsTable").Policy(Actions.DynamoDBRead),
        sls.Resource("tasksTable").Policy(Actions.DynamoDBCrud),
      ],
      events: [HttpApi.Delete("/lab/{labid}/task/{taskid}")],
    },

    getLabTasksFunction: {
      handler: handler("get-lab-tasks"),
      iamRoleStatements: [
        sls.Resource("tasksTable").Policy(Actions.DynamoDBRead),
      ],
      events: [HttpApi.Get("/lab/{labid}/tasks")],
    },

    getLabsFunction: {
      handler: handler("get-labs"),
      iamRoleStatements: [
        sls.Resource("labsTable").Policy(Actions.DynamoDBRead),
      ],
      events: [HttpApi.Get("/labs")],
    },

    pollTaskTestFunction: {
      handler: handler("poll-task-test"),
      iamRoleStatements: [
        sls.Resource("testRunsTable").Policy(Actions.DynamoDBRead),
      ],
      events: [HttpApi.Get("/lab/{labid}/task/{taskid}/test")],
    },

    runTaskTestFunction: {
      handler: handler("run-task-test"),
      iamRoleStatements: [
        sls.Resource("filesTable").Policy(Actions.DynamoDBCrud),
        sls.Resource("testRunsTable").Policy(Actions.DynamoDBCrud),
        sls.Resources("filesBucket", "filesBucket/*").Policy(Actions.S3Crud),
      ],
      events: [HttpApi.Post("/lab/{labid}/task/{taskid}/test")],
    },
  };
};
