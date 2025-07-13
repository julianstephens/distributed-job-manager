import { QueryClient } from "@tanstack/react-query";

export const convertUnixToDate = (
  unixTimestamp: number | undefined
): string | null => {
  if (!unixTimestamp) return null;
  const date = new Date(unixTimestamp * 1000);
  return date.toLocaleString();
};

export const TABLE_PAGE_SIZE = 10;

export const TaskStatus = {
  0: "Pending",
  1: "In Progress",
  2: "Completed",
  3: "Failed",
  4: "Cancelled",
} as const;

export const TaskRecurrence = {
  0: "Once",
  1: "Daily",
  2: "Weekly",
  3: "Monthly",
} as const;

export const sleep = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};

export const queryClient = new QueryClient();