import type { HTTPResponse, Job, JobRequest, TokenProp } from "@/lib/types";
import { $api } from "./client";

export const getJobs = async ({ token }: TokenProp): Promise<Job[]> => {
  if (!token) throw new Error("auth token missing");
  const res = await $api.get<HTTPResponse<Job[]>>("/jobs", {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const getJob = async ({
  token,
  id,
}: TokenProp & { id: string }): Promise<Job> => {
  const res = await $api.get<HTTPResponse<Job>>(`/jobs/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const createJob = async ({
  token,
  req,
}: TokenProp & { req: JobRequest }): Promise<Job> => {
  const res = await $api.post<HTTPResponse<Job>>("/jobs", req, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const deleteJob = async ({
  token,
  id,
}: TokenProp & { id: string }): Promise<string> => {
  const res = await $api.delete<HTTPResponse<string>>(`/jobs/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};
