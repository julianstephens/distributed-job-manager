import type {
  HTTPResponse,
  Job,
  JobCreateRequest,
  JobUpdateRequest,
  TokenProp,
} from "@/lib/types";
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
  if (!token) throw new Error("auth token missing");
  const res = await $api.get<HTTPResponse<Job>>(`/jobs/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const createJob = async ({
  token,
  req,
}: TokenProp & { req: JobCreateRequest }): Promise<Job> => {
  if (!token) throw new Error("auth token missing");
  const res = await $api.post<HTTPResponse<Job>>("/jobs", req, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const updateJob = async ({
  token,
  id,
  req,
}: TokenProp & { id: string; req: JobUpdateRequest }) => {
  if (!token) throw new Error("auth token missing");
  const res = await $api.patch<HTTPResponse<Job>>(`/jobs/${id}`, req, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};

export const deleteJob = async ({
  token,
  id,
}: TokenProp & { id: string }): Promise<string> => {
  if (!token) throw new Error("auth token missing");
  const res = await $api.delete<HTTPResponse<string>>(`/jobs/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data;
};
