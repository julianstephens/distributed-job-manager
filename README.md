# distributed-task-scheduler

```
{
   "title":  "Download file 1",
   "job": "curl https://dl.com/file1"
}

Runtime: 60min


{
   "title":  "Download file 2",
   "job": "curl https://dl.com/file2"
}

Runtime: 2min


```

Job Worker : lambda

/job_worker/submit => { job_id: "123" }

Queue: - Task 1 - Task 2 - Task 3

Worker 1 {job_id: "123", job: "curl https://dl.com/file2"}
Worker 2

---

Worker 3
...
Worker N
