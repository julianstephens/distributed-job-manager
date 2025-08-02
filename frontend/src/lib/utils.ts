import { QueryClient } from "@tanstack/react-query";

export const convertUnixToDate = (
  unixTimestamp: number | undefined
): string | null => {
  if (!unixTimestamp) return null;
  const date = new Date(unixTimestamp * 1000);
  return date.toLocaleString();
};

export const TABLE_PAGE_SIZE = 10;

export const JobStatus = {
  pending: "Pending",
  scheduled: "Scheduled",
  "in-progress": "In Progress",
  completed: "Completed",
  failed: "Failed",
  cancelled: "Cancelled",
} as const;

export const JobFrequency = {
  Once: "one-time",
  Daily: "daily",
  Weekly: "weekly",
  Monthly: "monthly",
} as const;

export const sleep = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};

export const queryClient = new QueryClient();

export const getKeyByValue = <T extends Record<string, any>>(
  object: T,
  value: T[keyof T]
): keyof T | undefined => {
  return (Object.keys(object) as Array<keyof T>).find(
    (key) => object[key] === value
  );
};

export const getJobStatusColor = (status: keyof typeof JobStatus) => {
  switch (status) {
    case "pending":
      return "gray";
    case "in-progress":
      return "blue";
    case "completed":
      return "green";
    case "failed":
      return "red";
    default:
      return "gray";
  }
};

export const displayDate = (date: string) => new Date(date).toLocaleString();

export const formatPayload = (payload: string) =>
  "```go\n" + payload.trimEnd() + "\n```";
