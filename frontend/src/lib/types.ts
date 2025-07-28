export interface ChildrenProps {
  children: React.ReactNode;
}

export interface HTTPResponse<T> {
  message: string;
  data: T;
}

export interface Job {
  job_id: string;
  user_id: string;
  job_name: string;
  frequency: string;
  status: string;
  payload: string;
  retry_count: number;
  max_retries: number;
  execution_time: string;
}

export type JobRequest = Omit<Job, "job_id" | "retry_count" | "status">;
