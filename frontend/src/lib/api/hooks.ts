import { useQuery } from "@tanstack/react-query";
import { getJob, getJobs } from "./queries";

export const useJobs = (userId?: string) => {
  return useQuery({
    queryKey: ["jobs", userId],
    queryFn: () => getJobs(userId!),
    enabled: !!userId,
  });
};

export const useJob = (jobId?: string) => {
  return useQuery({
    queryKey: ["job", jobId],
    queryFn: () => getJob(jobId!),
    enabled: !!jobId,
  });
};
