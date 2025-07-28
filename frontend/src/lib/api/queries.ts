import type { HTTPResponse, Job, JobRequest } from "@/lib/types";
import { $api } from "./client";

export const getJobs = async (userId: string): Promise<Job[]> => {
  const res = await $api.get<HTTPResponse<Job[]>>("/jobs", {
    params: {
      user_id: userId,
    },
  });
  return res.data.data;
};

export const getJob = async (id: string): Promise<Job> => {
  const res = await $api.get<HTTPResponse<Job>>(`/jobs/${id}`);
  return res.data.data;
};

export const createJob = async (req: JobRequest): Promise<Job> => {
  const res = await $api.post<HTTPResponse<Job>>("/jobs", req);
  return res.data.data;
};

export const deleteJob = async (id: string): Promise<string> => {
  const res = await $api.delete<HTTPResponse<string>>(`/jobs/${id}`);
  return res.data.data;
};
