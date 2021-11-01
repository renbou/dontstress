// Table "Labs"
interface LabsEntry {
  // Simple primary key
  id: string;
  // Attributes
  name: string;
}

// Table "Tasks"
interface TasksEntry {
  // Composite primary key
  labid: string;
  num: number;
  // Attributes
  name: string;
}

enum Languages {
  GPP = "G++",
  GCC = "GCC",
  JAVA = "Java",
  GO = "Go",
  PYTHON = "Python",
}

// Table "Files"
interface FilesEntry {
  // Simple primary key
  id: string;
  // Attributes
  lang: Languages;
}

enum TaskStatus {
  QUEUE = "QUEUE",
  RUNNING = "RUNNING",
  DONE = "DONE",
}

enum TestResult {
  OK = "OK",
  WA = "WA",
  ML = "ML",
  RE = "RE",
  CE = "CE",
}

// Table "TestRuns"
interface TestRunsEntry {
  // Simple primary key
  id: string;
  // Attributes
  status: TaskStatus;
  tests: {
    result: TestResult;
    info:
      | string
      | {
          test: string;
          expected: string;
          got: string;
        };
  }[];
}
