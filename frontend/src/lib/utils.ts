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
