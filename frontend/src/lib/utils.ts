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
