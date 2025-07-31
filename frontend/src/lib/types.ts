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
  job_description: string;
  frequency: string;
  status: string;
  payload: string;
  retry_count: number;
  max_retries: number;
  execution_time: string;
  created_at: string;
  updated_at: string;
}

export type JobRequest = Omit<
  Job,
  "job_id" | "user_id" | "retry_count" | "status" | "created_at" | "updated_at"
>;

export type TokenProp = {
  token: string;
};
